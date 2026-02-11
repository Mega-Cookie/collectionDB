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
	ID    int    `json:"id"`
	Title string `json:"title"`
	Plot  string `json:"plot"`
	Year  int    `json:"year"`
	Type  struct {
		ID   int    `json:"typeid"`
		Name string `json:"typename"`
	}
	Collection struct {
		ID   int    `json:"collid"`
		Name string `json:"collname"`
	}
	Genre struct {
		ID   int    `json:"genreid"`
		Name string `json:"genrename"`
	}
	IsDigital  bool      `json:"is_digital"`
	IsBooklet  bool      `json:"is_booklet"`
	MediaCount int       `json:"media_count"`
	Released   time.Time `json:"release_date"`
	Comment    string    `json:"comment"`
	AudioLangs string    `json:"audio_langs"`
	SubLangs   string    `json:"sub_langs"`
	RegionCode string    `json:"region_code"`
	BarCode    string    `json:"bar_code"`
	Imdb       string    `json:"imdbid"`
	CreatedAt  time.Time `json:"created_at"`
	EditedAt   time.Time `json:"edited_at"`
}

func ListEntries(db *sql.DB) (entries []Entry) {
	rows, err := db.Query("SELECT e.*, c.NAME AS COLLNAME, g.NAME AS GENRENAME, t.NAME AS TYPENAME FROM entries e LEFT OUTER JOIN genres g ON e.genreID = g.genreID LEFT OUTER JOIN collections c ON e.collectionID = c.collectionID LEFT OUTER JOIN mediatypes t ON e.mediatypeID = t.mediatypeID GROUP BY e.entryID")
	if err != nil {
		fmt.Println("error: Failed to retrieve entries")
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var entry Entry
		rows.Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.Comment, &entry.AudioLangs, &entry.SubLangs, &entry.IsDigital, &entry.Released, &entry.Collection.ID, &entry.Genre.ID, &entry.Type.ID, &entry.CreatedAt, &entry.EditedAt, &entry.Collection.Name, &entry.Genre.Name, &entry.Type.Name)
		entries = append(entries, entry)
	}
	return
}
func ShowCreateEntryPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		collections := collect.ListCollections(db)
		mediatypes := stockdata.ListMediatypes(db)
		genres := stockdata.ListGenres(db)
		c.Header("Cache-Control", "no-store")
		c.HTML(http.StatusOK, "entries/create.html", gin.H{
			"Collections": collections,
			"Genres":      genres,
			"Mediatypes":  mediatypes,
		})
	}
}
func CreateEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		plot := c.PostForm("plot")
		typeid := c.PostForm("typeid")
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

		_, err := db.Exec(`INSERT IGNORE INTO imdb (imdbID) VALUES (?)`, imdbid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to POST imdbID"})
			fmt.Println(err)
			return
		}

		_, err = db.Exec(`INSERT INTO entries (TITLE, YEAR, PLOT, mediatypeID, collectionID, genreID, IS_DIGITAL, IS_BOOKLET, MEDIARELEASEDATE, COMMENT, REGIONCODE, BARCODE, imdbID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, title, year, plot, typeid, collid, genreid, isdigital, isbooklet, released, comment, regioncode, barcode, imdbid)
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
		query := "SELECT e.*, c.NAME AS COLLNAME, g.NAME AS GENRENAME, t.NAME AS TYPENAME FROM entries e LEFT OUTER JOIN genres g ON e.genreID = g.genreID LEFT OUTER JOIN collections c ON e.collectionID = c.collectionID LEFT OUTER JOIN mediatypes t ON e.mediatypeID = t.mediatypeID WHERE e.entryID = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.IsDigital, &entry.Collection.ID, &entry.Genre.ID, &entry.Type.ID, &entry.CreatedAt, &entry.EditedAt, &entry.Collection.Name, &entry.Genre.Name, &entry.Type.Name)
		if err != nil {
			c.HTML(http.StatusNotFound, "entries/404.html", nil)
			fmt.Println(err)
			return
		}
		collections := collect.ListCollections(db)
		genres := stockdata.ListGenres(db)
		mediatypes := stockdata.ListMediatypes(db)
		entry.CreatedAt = small.SetTime(db, &entry.CreatedAt)
		entry.EditedAt = small.SetTime(db, &entry.EditedAt)
		c.Header("Cache-Control", "no-store")
		c.HTML(http.StatusOK, "entries/edit.html", gin.H{
			"Entry":       entry,
			"Collections": collections,
			"Genres":      genres,
			"Mediatypes":  mediatypes,
		})
	}
}
func EditEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		plot := c.PostForm("plot")
		typeid := c.PostForm("typeid")
		genreid := c.PostForm("genreid")
		year := c.PostForm("year")
		collid := c.PostForm("collid")
		isdigital := c.PostForm("is_digital") == "on"
		id := c.Param("id")
		updateTableQuery := `UPDATE entries SET TITLE = ?, YEAR = ?, PLOT = ?, mediatypeID = ?, genreID = ?, IS_DIGITAL = ?, collectionID = ?, EDITED_AT = CURRENT_TIMESTAMP where entryID = ?`
		_, err := db.Exec(updateTableQuery, title, year, plot, typeid, genreid, isdigital, collid, id)
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
		query := "SELECT e.*, c.NAME AS COLLNAME, g.NAME AS GENRENAME, t.NAME AS TYPENAME FROM entries e LEFT OUTER JOIN genres g ON e.genreID = g.genreID LEFT OUTER JOIN collections c ON e.collectionID = c.collectionID LEFT OUTER JOIN mediatypes t ON e.mediatypeID = t.mediatypeID WHERE e.entryID = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.IsDigital, &entry.Collection.ID, &entry.Genre.ID, &entry.Type.ID, &entry.CreatedAt, &entry.EditedAt, &entry.Collection.Name, &entry.Genre.Name, &entry.Type.Name)
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
