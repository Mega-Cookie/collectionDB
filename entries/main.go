package entries

import (
	"collectionDB/collect"
	"collectionDB/small"
	"collectionDB/stockdata"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID        int    `json:"entryID"`
	Title     string `json:"title"`
	Plot      string `json:"plot"`
	Year      int    `json:"year"`
	MediaType struct {
		ID   int    `json:"mediatypeid"`
		Name string `json:"mediatypename"`
	}
	CaseType struct {
		ID   int    `json:"casetypeid"`
		Name string `json:"casetypename"`
	}
	Collection struct {
		ID   int    `json:"collid"`
		Name string `json:"collname"`
	}
	Genre struct {
		ID   int    `json:"genreid"`
		Name string `json:"genrename"`
	}
	Publisher struct {
		ID   int    `json:"publisherid"`
		Name string `json:"publishername"`
	}
	IsDigital  bool      `json:"is_digital"`
	IsBooklet  bool      `json:"is_booklet"`
	MediaCount int       `json:"media_count"`
	Released   string    `json:"release_date"`
	Comment    string    `json:"comment"`
	AudioLangs string    `json:"audio_langs"`
	SubLangs   string    `json:"sub_langs"`
	RegionCode string    `json:"region_code"`
	BarCode    string    `json:"bar_code"`
	ImdbID     string    `json:"imdbid"`
	CreatedAt  time.Time `json:"created_at"`
	EditedAt   time.Time `json:"edited_at"`
}

func ListEntries(db *sql.DB) (entries []Entry) {
	querry := `SELECT e.*, c.NAME AS COLLNAME, g.NAME AS GENRENAME, mt.NAME AS MEDIATYPENAME, ct.NAME AS CASETYPENAME, p.NAME AS PUBNAME
				FROM entries e
				LEFT OUTER JOIN collections c ON e.collectionID = c.collectionID
				LEFT OUTER JOIN genres g ON e.genreID = g.genreID  
				LEFT OUTER JOIN mediatypes mt ON e.mediatypeID = mt.mediatypeID 
				LEFT OUTER JOIN casetypes ct on e.casetypeID = ct.casetypeID
				LEFT OUTER JOIN publishers p ON e.publisherID = p.publisherID 
				GROUP BY e.entryID`
	rows, err := db.Query(querry)
	if err != nil {
		fmt.Println("error: Failed to retrieve entries")
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var entry Entry
		rows.Scan(
			&entry.ID,
			&entry.Title,
			&entry.Year,
			&entry.Plot,
			&entry.Comment,
			&entry.AudioLangs,
			&entry.SubLangs,
			&entry.Released,
			&entry.MediaCount,
			&entry.IsDigital,
			&entry.IsBooklet,
			&entry.RegionCode,
			&entry.BarCode,
			&entry.Collection.ID,
			&entry.Genre.ID,
			&entry.MediaType.ID,
			&entry.CaseType.ID,
			&entry.Publisher.ID,
			&entry.ImdbID,
			&entry.CreatedAt,
			&entry.EditedAt,
			&entry.Collection.Name,
			&entry.Genre.Name,
			&entry.MediaType.Name,
			&entry.CaseType.Name,
			&entry.Publisher.Name)
		entries = append(entries, entry)
	}
	return
}
func ShowCreateEntryPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		collections := collect.ListCollections(db)
		mediatypes := stockdata.ListMediaTypes(db)
		casetypes := stockdata.ListCaseTypes(db)
		publishers := stockdata.ListPublishers(db)
		genres := stockdata.ListGenres(db)
		c.Header("Cache-Control", "no-store")
		c.HTML(http.StatusOK, "entries/create.html", gin.H{
			"Collections": collections,
			"Genres":      genres,
			"MediaTypes":  mediatypes,
			"CaseTypes":   casetypes,
			"Publishers":  publishers,
		})
	}
}
func CreateEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		plot := c.PostForm("plot")
		mediatypeid := c.PostForm("mediatypeid")
		casetypeid := c.PostForm("casetypeid")
		genreid := c.PostForm("genreid")
		year := c.PostForm("year")
		collid := c.PostForm("collid")
		isdigital := c.PostForm("is_digital") == "on"
		isbooklet := c.PostForm("is_booklet") == "on"
		released := c.PostForm("release_date")
		comment := c.PostForm("comment")
		regioncode := c.PostForm("region_code")
		barcode := c.PostForm("bar_code")
		imdbid := c.PostForm("imdbid")
		publisherid := c.PostForm("publisherid")
		if imdbid != "" {
			_, err := db.Exec(`INSERT OR IGNORE INTO imdb (imdbID) VALUES (?)`, imdbid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to POST imdbID"})
				fmt.Println(err)
				return
			}
		}
		_, err := db.Exec(`INSERT INTO entries (TITLE, YEAR, PLOT, mediatypeID, casetypeID, collectionID, genreID, IS_DIGITAL, IS_BOOKLET, MEDIARELEASEDATE, COMMENT, REGIONCODE, BARCODE, imdbID, publisherID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, title, year, plot, mediatypeid, casetypeid, collid, genreid, isdigital, isbooklet, released, comment, regioncode, barcode, imdbid, publisherid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/")
	}
}
func ShowEditEntryPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var entry Entry
		query := `SELECT e.*, c.NAME AS COLLNAME, g.NAME AS GENRENAME, mt.NAME AS MEDIATYPENAME, ct.NAME AS CASETYPENAME, p.NAME AS PUBNAME
				FROM entries e
				LEFT OUTER JOIN collections c ON e.collectionID = c.collectionID
				LEFT OUTER JOIN genres g ON e.genreID = g.genreID
				LEFT OUTER JOIN mediatypes mt ON e.mediatypeID = mt.mediatypeID
				LEFT OUTER JOIN casetypes ct ON e.casetypeID = ct.casetypeID
				LEFT OUTER JOIN publishers p ON e.publisherID = p.publisherID
				WHERE e.entryID = ?`
		err := db.QueryRow(query, id).Scan(
			&entry.ID,
			&entry.Title,
			&entry.Year,
			&entry.Plot,
			&entry.Comment,
			&entry.AudioLangs,
			&entry.SubLangs,
			&entry.Released,
			&entry.MediaCount,
			&entry.IsDigital,
			&entry.IsBooklet,
			&entry.RegionCode,
			&entry.BarCode,
			&entry.Collection.ID,
			&entry.Genre.ID,
			&entry.MediaType.ID,
			&entry.CaseType.ID,
			&entry.Publisher.ID,
			&entry.ImdbID,
			&entry.CreatedAt,
			&entry.EditedAt,
			&entry.Collection.Name,
			&entry.Genre.Name,
			&entry.MediaType.Name,
			&entry.CaseType.Name,
			&entry.Publisher.Name)
		if err != nil {
			c.HTML(http.StatusNotFound, "entries/404.html", nil)
			fmt.Println(err)
			return
		}
		collections := collect.ListCollections(db)
		genres := stockdata.ListGenres(db)
		mediatypes := stockdata.ListMediaTypes(db)
		entry.CreatedAt = small.SetTime(db, &entry.CreatedAt)
		entry.EditedAt = small.SetTime(db, &entry.EditedAt)
		c.Header("Cache-Control", "no-store")
		c.HTML(http.StatusOK, "entries/edit.html", gin.H{
			"Entry":       entry,
			"Collections": collections,
			"Genres":      genres,
			"MediaTypes":  mediatypes,
		})
	}
}
func EditEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		plot := c.PostForm("plot")
		mediatypeid := c.PostForm("mediatypeid")
		genreid := c.PostForm("genreid")
		year := c.PostForm("year")
		collid := c.PostForm("collid")
		isdigital := c.PostForm("is_digital") == "on"
		id := c.Param("id")
		updateTableQuery := `UPDATE entries SET TITLE = ?, YEAR = ?, PLOT = ?, mediatypeID = ?, genreID = ?, IS_DIGITAL = ?, collectionID = ?, EDITED_AT = CURRENT_TIMESTAMP where entryID = ?`
		_, err := db.Exec(updateTableQuery, title, year, plot, mediatypeid, genreid, isdigital, collid, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit entry"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/")
	}
}
func ViewEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var entry Entry
		query := `SELECT e.*, c.NAME AS COLLNAME, g.NAME AS GENRENAME, mt.NAME AS MEDIATYPENAME, ct.NAME AS CASETYPENAME, p.NAME AS PUBNAME
				FROM entries e
				LEFT OUTER JOIN collections c ON e.collectionID = c.collectionID
				LEFT OUTER JOIN genres g ON e.genreID = g.genreID
				LEFT OUTER JOIN mediatypes mt ON e.mediatypeID = mt.mediatypeID
				LEFT OUTER JOIN casetypes ct ON e.casetypeID = ct.casetypeID
				LEFT OUTER JOIN publishers p ON e.publisherID = p.publisherID
				WHERE e.entryID = ?`
		err := db.QueryRow(query, id).Scan(
			&entry.ID,
			&entry.Title,
			&entry.Year,
			&entry.Plot,
			&entry.Comment,
			&entry.AudioLangs,
			&entry.SubLangs,
			&entry.Released,
			&entry.MediaCount,
			&entry.IsDigital,
			&entry.IsBooklet,
			&entry.RegionCode,
			&entry.BarCode,
			&entry.Collection.ID,
			&entry.Genre.ID,
			&entry.MediaType.ID,
			&entry.CaseType.ID,
			&entry.Publisher.ID,
			&entry.ImdbID,
			&entry.CreatedAt,
			&entry.EditedAt,
			&entry.Collection.Name,
			&entry.Genre.Name,
			&entry.MediaType.Name,
			&entry.CaseType.Name,
			&entry.Publisher.Name)
		if err != nil {
			c.HTML(http.StatusNotFound, "entries/404.html", nil)
			fmt.Println(err)
			return
		}
		entry.CreatedAt = small.SetTime(db, &entry.CreatedAt)
		entry.EditedAt = small.SetTime(db, &entry.EditedAt)
		c.Header("Cache-Control", "no-store")
		c.HTML(http.StatusOK, "entries/view.html", entry)
	}
}
func DeleteEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM entries WHERE entryID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entry"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/")
	}
}
