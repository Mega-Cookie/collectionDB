package entries

import (
	"collectionDB/collect"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Plot      string    `json:"plot"`
	Medium    string    `json:"medium"`
	Year      int       `json:"year"`
	CollID    string    `json:"collid"`
	CollName  string    `json:"collname"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
	IsDigital bool      `json:"is_digital"`
}
type Collection struct {
	CollID      int       `json:"collid"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	EditedAt    time.Time `json:"edited_at"`
}

func ShowCreateEntryPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_entry.html", gin.H{})
}
func CreateEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		plot := c.PostForm("plot")
		medium := c.PostForm("medium")
		year := c.PostForm("year")
		collid := c.PostForm("collid")
		isDigital := c.PostForm("is_digital") == "on"
		_, err := db.Exec(`INSERT INTO entries (TITLE, YEAR, PLOT, MEDIUM, IS_DIGITAL, collectionID) VALUES (?, ?, ?, ?, ?, ?)`, title, year, plot, medium, isDigital, collid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func EditEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		plot := c.PostForm("plot")
		media := c.PostForm("media")
		year := c.PostForm("year")
		isDigital := c.PostForm("is_digital") == "on"
		id := c.Param("id")
		updateTableQuery := `UPDATE entries SET title = ?, year = ?, content = ?, media = ?, is_digital = ?, edited_at = CURRENT_TIMESTAMP where id = ?`
		_, err := db.Exec(updateTableQuery, title, year, plot, media, isDigital, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit entry"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func List(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		collections := collect.ListCollections(db)
		rows, err := db.Query("SELECT e.*, c.NAME AS COLLNAME FROM `entries` e JOIN collections c on c.collectionID = e.collectionID")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve entries"})
			return
		}
		defer rows.Close()
		var entries []Entry
		for rows.Next() {
			var entry Entry
			if err := rows.Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.Medium, &entry.IsDigital, &entry.CollID, &entry.CreatedAt, &entry.EditedAt, &entry.CollName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while scanning entries"})
				return
			}
			entries = append(entries, entry)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Entries":     entries,
			"Collections": collections,
		})
	}
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
		query := "SELECT entryID, TITLE, YEAR, PLOT, MEDIUM, collectionID, CREATED_AT, EDITED_AT FROM entries WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.Medium, &entry.CollName, &entry.CreatedAt, &entry.EditedAt)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		c.HTML(http.StatusOK, "preview.html", entry)
	}
}
func ShowEditEntryPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var entry Entry
		query := "SELECT e.*, c.NAME AS COLLNAME FROM `entries` e JOIN collections c on c.collectionID = e.collectionID WHERE entryID = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Plot, &entry.Medium, &entry.IsDigital, &entry.CollID, &entry.CreatedAt, &entry.EditedAt, &entry.CollName)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		defer collect.ListCollections(db)
		c.HTML(http.StatusOK, "edit_entry.html", gin.H{
			"Entry": entry,
		})
	}
}
