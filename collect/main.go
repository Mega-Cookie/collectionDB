package collect

import (
	"collectionDB/small"
	"collectionDB/stockdata"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Collection struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Entrycount  int    `json:"entrycount"`
	Category    struct {
		ID   int    `json:"catid"`
		Name string `json:"catname"`
	}
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
}

func ListCollections(db *sql.DB) (collections []Collection) {
	rows, err := db.Query("SELECT c.*, count(e.collectionID) AS ENTRYCOUNT, ca.NAME AS CATNAME FROM collections c LEFT OUTER JOIN categories ca ON ca.categoryID = c.categoryID LEFT OUTER JOIN entries e ON c.collectionID = e.collectionID GROUP BY c.collectionID")
	if err != nil {
		fmt.Println("error: Failed to retrieve collections")
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var collection Collection
		rows.Scan(&collection.ID, &collection.Name, &collection.Description, &collection.Category.ID, &collection.CreatedAt, &collection.EditedAt, &collection.Entrycount, &collection.Category.Name)
		collections = append(collections, collection)
	}
	return
}
func ShowCreateCollectionPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories := stockdata.ListCategories(db)
		c.HTML(http.StatusOK, "collections/create.html", gin.H{
			"Categories": categories,
		})
	}
}
func CreateCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		catid := c.PostForm("catid")
		description := c.PostForm("description")
		_, err := db.Exec(`INSERT INTO collections (NAME, categoryID, DESCRIPTION) VALUES (?, ?, ?)`, name, catid, description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection!"})
			fmt.Println(err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func ShowEditCollectionPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories := stockdata.ListCategories(db)
		id := c.Param("id")
		var collection Collection
		query := "SELECT c.*, count(e.collectionID) AS ENTRYCOUNT, ca.NAME AS CATNAME FROM collections c LEFT OUTER JOIN categories ca ON ca.categoryID = c.categoryID LEFT OUTER JOIN entries e ON c.collectionID = e.collectionID WHERE c.collectionID = ?"
		err := db.QueryRow(query, id).Scan(&collection.ID, &collection.Name, &collection.Description, &collection.Category.ID, &collection.CreatedAt, &collection.EditedAt, &collection.Entrycount, &collection.Category.Name)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			fmt.Println(err)
			return
		}
		collection.CreatedAt = small.SetTime(db, &collection.CreatedAt)
		collection.EditedAt = small.SetTime(db, &collection.EditedAt)
		c.HTML(http.StatusOK, "collections/edit.html", gin.H{
			"Collection": collection,
			"Categories": categories,
		})
	}
}
func EditCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		catid := c.PostForm("catid")
		description := c.PostForm("description")
		id := c.Param("id")
		updateTableQuery := `UPDATE collections SET NAME = ?, categoryID = ?, DESCRIPTION = ?, EDITED_AT = CURRENT_TIMESTAMP where collectionID = ?`
		_, err := db.Exec(updateTableQuery, name, catid, description, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit entry"})
			fmt.Println(err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
func ViewCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var collection Collection
		query := "SELECT c.*, ca.NAME AS CATNAME FROM collections c LEFT OUTER JOIN categories ca ON ca.categoryID = c.categoryID WHERE c.collectionID = ?"
		err := db.QueryRow(query, id).Scan(&collection.ID, &collection.Name, &collection.Description, &collection.Category.ID, &collection.CreatedAt, &collection.EditedAt, &collection.Category.Name)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			fmt.Println(err)
			return
		}
		collection.CreatedAt = small.SetTime(db, &collection.CreatedAt)
		collection.EditedAt = small.SetTime(db, &collection.EditedAt)
		c.HTML(http.StatusOK, "collections/view.html", collection)
	}
}
func DeleteCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM collections WHERE collectionID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection!"})
			fmt.Println(err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
