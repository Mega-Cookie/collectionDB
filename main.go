package main

import (
	"collectionDB/collect"
	"collectionDB/entries"
	"collectionDB/small"
	"database/sql"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func setStockData() {
	var err error
	_, err = db.Exec(`INSERT OR IGNORE INTO mediatypes (NAME)
		VALUES
		('CD'),
		('BlueRay'),
		('DVD'),
		('Manga'),
		('Comic'),
		('Book');`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`INSERT OR IGNORE INTO categories (NAME)
		VALUES
		('Movie'),
		('TV-Series'),
		('Music'),
		('Literature');`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`INSERT OR IGNORE INTO genres (NAME)
		VALUES
		('Fantasy'),
		('Romance'),
		('Action'),
		('Science Fiction'),
		('Musical'),
		('Horror');`)
	if err != nil {
		log.Fatal(err)
	}
}
func initDB(databasefile string) {
	var err error
	database, _ := filepath.Abs(databasefile)
	db, err = sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery := `CREATE TABLE IF NOT EXISTS mediatypes (
		typeID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TABLE IF NOT EXISTS genres (
		genreID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS categories (
		categoryID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS entries (
		entryID INTEGER PRIMARY KEY AUTOINCREMENT,
		TITLE TEXT NOT NULL,
		YEAR INTEGER NOT NULL,
		PLOT TEXT NOT NULL,
		IS_DIGITAL BOOL NOT NULL DEFAULT 0,
		collectionID INTEGER DEFAULT NULL,
		genreID INT DEFAULT NULL,
		typeID INTEGER DEFAULT NULL,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(collectionID) REFERENCES collections(collectionID),
		FOREIGN KEY(genreID) REFERENCES genres(genreID),
		FOREIGN KEY(typeID) REFERENCES mediatypes(typeID)
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS collections (
		collectionID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME TEXT UNIQUE,
		DESCRIPTION TEXT,
		categoryID INT DEFAULT NULL,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(categoryID) REFERENCES categories(categoryID)
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS info (
		instanceID INTEGER PRIMARY KEY,
		VERSION STRING,
		HOSTNAME STRING UNIQUE,
		OS STRING,
		ARCH STRING,
		GOVERSION STRING,
		TIMEZONE STRING
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
func ShowList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		collections := collect.ListCollections(db)
		entries := entries.ListEntries(db)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Entries":     entries,
			"Collections": collections,
		})
	}
}
func main() {
	config := small.Configure()
	initDB(config.Database)
	setStockData()
	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	small.SetSystemInfo(db)
	router := gin.Default()
	router.LoadHTMLGlob(config.Templates)
	router.GET("/", ShowList(db))
	router.GET("/create_entry", entries.ShowCreateEntryPage(db))
	router.POST("/create_entry", entries.CreateEntry(db))
	router.GET("/entries/:id/edit", entries.ShowEditEntryPage(db))
	router.POST("/entries/:id/edit", entries.EditEntry(db))
	router.POST("/entries/:id/delete", entries.DeleteEntry(db))
	router.GET("/entries/:id", entries.ViewEntry(db))
	router.GET("/create_collection", collect.ShowCreateCollectionPage(db))
	router.POST("/create_collection", collect.CreateCollection(db))
	router.GET("/collections/:id/edit", collect.ShowEditCollectionPage(db))
	router.POST("/collections/:id/edit", collect.EditCollection(db))
	router.POST("/collections/:id/delete", collect.DeleteCollection(db))
	router.GET("/collections/:id", collect.ViewCollection(db))
	log.Printf("Server is running on http://%s", config.Listen)
	log.Printf("Acessing SQLite: %s", config.Database)
	if err := router.Run(config.Listen); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
