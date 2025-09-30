package entries

import (
	"collectionDB/collect"
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
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
	IsDigital bool      `json:"is_digital"`
}

func ShowCreateEntryPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		collections := collect.ListCollections(db)
		c.HTML(http.StatusOK, "create_entry.html", gin.H{
			"Collections": collections,
		})
	}
}
func CreateEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		Title := c.PostForm("title")
		Plot := c.PostForm("plot")
		Typeid := c.PostForm("typeid")
		Genreid := c.PostForm("genreid")
		Year := c.PostForm("year")
		Collid := c.PostForm("collid")
		IsDigital := c.PostForm("is_digital") == "on"
		_, err := db.Exec(`INSERT INTO entries (TITLE, YEAR, PLOT, TypeID, GenreID, IS_DIGITAL, collectionID) VALUES (?, ?, ?, ?, ?, ?)`, Title, Year, Plot, Typeid, Genreid, IsDigital, Collid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func EditEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		Title := c.PostForm("title")
		Plot := c.PostForm("plot")
		Typeid := c.PostForm("typeid")
		Genreid := c.PostForm("genreid")
		Year := c.PostForm("year")
		Collid := c.PostForm("collid")
		IsDigital := c.PostForm("is_digital") == "on"
		id := c.Param("id")
		updateTableQuery := `UPDATE entries SET TITLE = ?, YEAR = ?, PLOT = ?, TypeID = ?, GenreID = ?, IS_DIGITAL = ?, collectionID = ?, , EDITED_AT = CURRENT_TIMESTAMP where entryID = ?`
		_, err := db.Exec(updateTableQuery, Title, Year, Plot, Typeid, Genreid, IsDigital, Collid, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit entry"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func ListEntries(db *sql.DB) (entries []Entry) {
	rows, err := db.Query("SELECT e.*, c.NAME AS COLLNAME FROM `entries` e JOIN collections c on c.collectionID = e.collectionID")
	if err != nil {
		fmt.Println("error: Failed to retrieve entries")
	}
	defer rows.Close()
	for rows.Next() {
		var entry Entry
		rows.Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.Type.ID, &entry.IsDigital, &entry.Collection.ID, entry.Genre.ID, entry.Type.ID, &entry.CreatedAt, &entry.EditedAt, &entry.Collection.Name, &entry.Genre.Name, &entry.Type.Name)
		entries = append(entries, entry)
	}
	return
}
func DeleteEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM entries WHERE entryID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entry"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func PreviewSharedEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var entry Entry
		query := "SELECT entryID, TITLE, YEAR, PLOT, IS_DIGITAL, collectionID GenreID, TypeID, CREATED_AT, EDITED_AT FROM entries WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.IsDigital, &entry.Collection.ID, entry.Genre.ID, entry.Type.ID, &entry.CreatedAt, &entry.EditedAt)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		c.HTML(http.StatusOK, "preview.html", entry)
	}
}
func ShowEditEntryPage(db *sql.DB) gin.HandlerFunc {
	collections := collect.ListCollections(db)
	genres := stockdata.ListGenres(db)
	mediatypes := stockdata.ListMediatypes(db)
	return func(c *gin.Context) {
		id := c.Param("id")
		var entry Entry
		query := "SELECT e.*, c.NAME AS COLLNAME FROM `entries` e JOIN collections c on c.collectionID = e.collectionID WHERE entryID = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.IsDigital, &entry.Collection.ID, entry.Genre.ID, entry.Type.ID, &entry.CreatedAt, &entry.EditedAt, &entry.Collection.Name, &entry.Genre.Name, &entry.Type.Name)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		c.HTML(http.StatusOK, "edit_entry.html", gin.H{
			"Entry":       entry,
			"Collections": collections,
			"Genres":      genres,
			"Mediatypes":  mediatypes,
		})
	}
}
