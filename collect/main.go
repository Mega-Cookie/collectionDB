package collect

import (
	"collectionDB/collect"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Collection struct {
	CollID      int       `json:"collid"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	EditedAt    time.Time `json:"edited_at"`
}

func ListCollections(db *sql.DB) (collections []Collection) {
	rows, err := db.Query("SELECT collectionID, NAME, TYPE, DESCRIPTION, CREATED_AT, EDITED_AT FROM collections")
	if err != nil {
		fmt.Println("error: Failed to retrieve collections")
	}
	defer rows.Close()
	for rows.Next() {
		var collection Collection
		rows.Scan(&collection.CollID, &collection.Name, &collection.Type, &collection.Description, &collection.CreatedAt, &collection.EditedAt)
		collections = append(collections, collection)
	}
	return
}

func ShowCreateCollectionPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_collection.html", nil)
}
func CreateCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		Name := c.PostForm("name")
		Type := c.PostForm("type")
		Description := c.PostForm("description")
		_, err := db.Exec(`INSERT INTO collections (NAME, TYPE, DESCRIPTION) VALUES (?, ?, ?)`, Name, Type, Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection!"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func ShowEditCollectionPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var entry Collection
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
func EditCollection(db *sql.DB) gin.HandlerFunc {
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
func DeleteCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM collections WHERE entryID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection!"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
