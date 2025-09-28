package collections

import (
	"database/sql"
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

func ListCollections(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT collectionID, NAME, TYPE, CREATED_AT, EDITED_AT FROM collections")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collections"})
			return
		}
		defer rows.Close()
		var collections []Collection
		for rows.Next() {
			var collection Collection
			if err := rows.Scan(&collection.CollID, &collection.Name, &collection.Type, &collection.CreatedAt, &collection.EditedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while scanning collections"})
				return
			}
			collections = append(collections, collection)
		}
		c.HTML(http.StatusOK, "collections.html", gin.H{
			"Collections": collections,
		})
	}
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
