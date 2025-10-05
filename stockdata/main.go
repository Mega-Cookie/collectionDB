package stockdata

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Mediatypes struct {
	ID        int       `json:"typeid"`
	Name      string    `json:"name"`
	IsStock   bool      `json:"is_stock"`
	CreatedAt time.Time `json:"created_at"`
}
type Categories struct {
	ID        int       `json:"catid"`
	Name      string    `json:"name"`
	IsStock   bool      `json:"is_stock"`
	CreatedAt time.Time `json:"created_at"`
}
type Genres struct {
	ID        int       `json:"genreid"`
	Name      string    `json:"name"`
	IsStock   bool      `json:"is_stock"`
	CreatedAt time.Time `json:"created_at"`
}

func ListMediatypes(db *sql.DB) (mediatypes []Mediatypes) {
	rows, err := db.Query("SELECT * FROM mediatypes")
	if err != nil {
		fmt.Println("error: Failed to retrieve collections")
	}
	defer rows.Close()
	for rows.Next() {
		var mediatype Mediatypes
		rows.Scan(&mediatype.ID, &mediatype.Name, &mediatype.IsStock, &mediatype.CreatedAt)
		mediatypes = append(mediatypes, mediatype)
	}
	fmt.Println(mediatypes)
	return
}
func ListCategories(db *sql.DB) (categories []Categories) {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		fmt.Println("error: Failed to retrieve categories")
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var category Categories
		rows.Scan(&category.ID, &category.Name, &category.IsStock, &category.CreatedAt)
		categories = append(categories, category)
	}
	return
}
func ListGenres(db *sql.DB) (genres []Genres) {
	rows, err := db.Query("SELECT * FROM genres")
	if err != nil {
		fmt.Println("error: Failed to retrieve genres")
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var genre Genres
		rows.Scan(&genre.ID, &genre.Name, &genre.IsStock, &genre.CreatedAt)
		genres = append(genres, genre)
	}
	return
}
func CreateType(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		_, err := db.Exec(`INSERT INTO mediatypes (NAME) VALUES (?)`, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mediatype"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/stock")
	}
}
func DeleteType(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM mediatypes WHERE typeID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mediatype"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/stock")
	}
}
func CreateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		_, err := db.Exec(`INSERT INTO categories (NAME) VALUES (?)`, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/stock")
	}
}
func DeleteCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM categories WHERE categoryID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/stock")
	}
}
func CreateGenre(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		_, err := db.Exec(`INSERT INTO genres (NAME) VALUES (?)`, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/stock")
	}
}
func DeleteGenre(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM genres WHERE genreID = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete genre"})
			fmt.Println(err)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, "/stock")
	}
}
