package stockdata

import (
	"database/sql"
	"fmt"
)

type Mediatypes struct {
	ID   int    `json:"typeid"`
	Name string `json:"name"`
}

func ListMediatypes(db *sql.DB) (mediatypes []Mediatypes) {
	rows, err := db.Query("SELECT * FROM mediatypes")
	if err != nil {
		fmt.Println("error: Failed to retrieve collections")
	}
	defer rows.Close()
	for rows.Next() {
		var mediatype Mediatypes
		rows.Scan(&mediatype.ID, &mediatype.Name)
		mediatypes = append(mediatypes, mediatype)
	}
	return
}

type Categories struct {
	ID   int    `json:"typeid"`
	Name string `json:"name"`
}

func ListCategories(db *sql.DB) (categories []Categories) {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		fmt.Println("error: Failed to retrieve categories")
	}
	defer rows.Close()
	for rows.Next() {
		var category Categories
		rows.Scan(&category.ID, &category.Name)
		categories = append(categories, category)
	}
	return
}

type Genres struct {
	ID   int    `json:"genreid"`
	Name string `json:"name"`
}

func ListGenres(db *sql.DB) (genres []Genres) {
	rows, err := db.Query("SELECT * FROM genres")
	if err != nil {
		fmt.Println("error: Failed to retrieve genres")
	}
	defer rows.Close()
	for rows.Next() {
		var genre Genres
		rows.Scan(&genre.ID, &genre.Name)
		genres = append(genres, genre)
	}
	return
}
