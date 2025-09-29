package collect

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Collection struct {
	CollID      int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Entrycount  int       `json:"entrycount"`
	CreatedAt   time.Time `json:"created_at"`
	EditedAt    time.Time `json:"edited_at"`
}

func ListCollections(db *sql.DB) (collections []Collection) {
	rows, err := db.Query("SELECT c.*, count(e.collectionID) AS ENTRYCOUNT FROM `collections` c LEFT OUTER JOIN entries e on c.collectionID = e.collectionID GROUP BY c.collectionID")
	if err != nil {
		fmt.Println("error: Failed to retrieve collections")
	}
	defer rows.Close()
	for rows.Next() {
		var collection Collection
		rows.Scan(&collection.CollID, &collection.Name, &collection.Type, &collection.Description, &collection.CreatedAt, &collection.EditedAt, &collection.Entrycount)
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
		var collection Collection
		query := "SELECT c.*, count(e.collectionID) AS ENTRYCOUNT FROM `collections` c LEFT OUTER JOIN entries e on c.collectionID = e.collectionID WHERE c.collectionID = ?"
		err := db.QueryRow(query, id).Scan(&collection.CollID, &collection.Name, &collection.Type, &collection.Description, &collection.CreatedAt, &collection.EditedAt, &collection.Entrycount)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		c.HTML(http.StatusOK, "edit_collection.html", gin.H{
			"Collection": collection,
		})
	}
}
func EditCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		Name := c.PostForm("name")
		Type := c.PostForm("type")
		Description := c.PostForm("description")
		id := c.Param("id")
		updateTableQuery := `UPDATE collections SET NAME = ?, TYPE = ?, Description = ?, EDITED_AT = CURRENT_TIMESTAMP where collectionID = ?`
		_, err := db.Exec(updateTableQuery, Name, Type, Description, id)
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
		_, err := db.Exec(`DELETE FROM collections WHERE collectionID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection!"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
