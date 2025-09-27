package entries

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Media     string    `json:"media"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
	IsDigital bool      `json:"is_digital"`
}

func ShowCreateEntryPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}
func CreateEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		media := c.PostForm("media")
		year := c.PostForm("year")
		isDigital := c.PostForm("is_digital") == "on"
		_, err := db.Exec(`INSERT INTO entries (title, year, content, media, is_digital) VALUES (?, ?, ?, ?, ?)`, title, year, content, media, isDigital)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
			return
		}
		c.Redirect(http.StatusOK, "/")
	}
}
func EditEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		media := c.PostForm("media")
		year := c.PostForm("year")
		isDigital := c.PostForm("is_digital") == "on"
		id := c.Param("id")
		updateTableQuery := `UPDATE entries SET title = ?, year = ?, content = ?, media = ?, is_digital = ?, edited_at = CURRENT_TIMESTAMP where id = ?`
		_, err := db.Exec(updateTableQuery, title, year, content, media, isDigital, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit entry"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func ListEntries(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, year, content, media, is_digital, created_at FROM entries")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve entries"})
			return
		}
		defer rows.Close()
		var entries []Entry
		for rows.Next() {
			var entry Entry
			if err := rows.Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Content, &entry.Media, &entry.IsDigital, &entry.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while scanning entries"})
				return
			}
			entries = append(entries, entry)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Entries": entries,
		})
	}
}
func DeleteEntry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM entries WHERE id = ?`, id)
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
		query := "SELECT id, title, year, content, media, created_at, edited_at FROM entries WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Content, &entry.Media, &entry.CreatedAt, &entry.EditedAt)
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
		query := "SELECT id, title, year, content, is_digital, media, created_at, edited_at FROM entries WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Year, &entry.Content, &entry.IsDigital, &entry.Media, &entry.CreatedAt, &entry.EditedAt)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		c.HTML(http.StatusOK, "edit.html", entry)
	}
}
