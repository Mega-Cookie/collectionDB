package main

import (
	"collectionDB/collect"
	"collectionDB/entries"
	"collectionDB/small"
	"collectionDB/stockdata"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func setStockData() {
	var err error
	_, err = db.Exec(`INSERT OR IGNORE INTO categories (NAME, STOCK)
		VALUES
		('Movie', '1'),
		('TV-Series', '1');`)
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
		mediatypeID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		STOCK BOOLEAN NOT NULL DEFAULT 0,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS casetypes (
		casetypeID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		STOCK BOOLEAN NOT NULL DEFAULT 0,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS genres (
		genreID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		STOCK BOOLEAN NOT NULL DEFAULT 0,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS categories (
		categoryID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		STOCK BOOLEAN NOT NULL DEFAULT 0,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS collections (
		collectionID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME TEXT UNIQUE,
		DESCRIPTION TEXT,
		categoryID INTEGER DEFAULT NULL,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(categoryID) REFERENCES categories(categoryID)
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS imdb (
		imdbID TEXT PRIMARY KEY,
		RATING FLOAT,
		TITLE STRING UNIQUE,
		YEAR INTEGER,
		TAGLINE TEXT,
		PLOT TEXT,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS publishers (
		publisherID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		STOCK BOOLEAN NOT NULL DEFAULT 0,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TRIGGER IF NOT EXISTS abort_delete_stockcat
		BEFORE DELETE ON categories
		WHEN OLD.STOCK = 1
		BEGIN
    		SELECT RAISE(ABORT, 'You can''t delete system stock data');
		END
		;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS entries (
		entryID INTEGER PRIMARY KEY AUTOINCREMENT,
		TITLE TEXT NOT NULL,
		YEAR INTEGER NOT NULL,
		PLOT TEXT NOT NULL,
		COMMENT TEXT,
		AUDIOLANGS TEXT,
		SUBTITLELANGS TEXT,
		RELEASED DATE NOT NULL,
		IS_DIGITAL BOOLEAN NOT NULL DEFAULT 0,
		collectionID INTEGER DEFAULT NULL,
		genreID INTEGER DEFAULT NULL,
		mediatypeID INTEGER DEFAULT NULL,
		MEDIACOUNT INTEGER NOT NULL DEFAULT 1,
		IS_BOOKLET BOOLEAN NOT NULL DEFAULT 0,
		casetypeID INTEGER DEFAULT NULL,
		publisherID INTEGER DEFAULT NULL,
		REGIONCODE INTEGER,
		BARCODE INTEGER,
		imdbID TEXT,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(collectionID) REFERENCES collections(collectionID),
		FOREIGN KEY(genreID) REFERENCES genres(genreID),
		FOREIGN KEY(mediatypeID) REFERENCES mediatypes(mediatypeID),
		FOREIGN KEY(casetypeID) REFERENCES casetypes(casetypeID),
		FOREIGN KEY(publisherID) REFERENCES publishers(publisherID),
		FOREIGN KEY(imdbID) REFERENCES imdb(imdbID)
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
		SQLITEVERSION STRING,
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
func ShowStockList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mediatypes := stockdata.ListMediatypes(db)
		categories := stockdata.ListCategories(db)
		genres := stockdata.ListGenres(db)
		c.HTML(http.StatusOK, "stock/index.html", gin.H{
			"Mediatypes": mediatypes,
			"Categories": categories,
			"Genres":     genres,
		})
	}
}
func ShowAbout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		systeminfo := small.GetSystemInfo(db)
		c.HTML(http.StatusOK, "about.html", gin.H{
			"Info": systeminfo,
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
	if config.IsReverseProxy {
		router.SetTrustedProxies([]string{"127.0.0.1"})
	}
	router.Static("/static", config.Static)
	templates := fmt.Sprintf("%s/templates/**/*.html", config.Static)
	router.LoadHTMLGlob(templates)
	router.GET("/", ShowList(db))
	router.GET("/stock", ShowStockList(db))
	router.GET("/about", ShowAbout(db))
	router.POST("/stock/mediatype/create", stockdata.CreateType(db))
	router.POST("/stock/category/create", stockdata.CreateCategory(db))
	router.POST("/stock/genre/create", stockdata.CreateGenre(db))
	router.POST("/stock/mediatype/:id/delete", stockdata.DeleteType(db))
	router.POST("/stock/category/:id/delete", stockdata.DeleteCategory(db))
	router.POST("/stock/genre/:id/delete", stockdata.DeleteGenre(db))
	router.GET("/entries/:id", entries.ViewEntry(db))
	router.GET("/create_entry", entries.ShowCreateEntryPage(db))
	router.POST("/create_entry", entries.CreateEntry(db))
	router.GET("/entries/:id/edit", entries.ShowEditEntryPage(db))
	router.POST("/entries/:id/edit", entries.EditEntry(db))
	router.POST("/entries/:id/delete", entries.DeleteEntry(db))
	router.GET("/collections/:id", collect.ViewCollection(db))
	router.GET("/create_collection", collect.ShowCreateCollectionPage(db))
	router.POST("/create_collection", collect.CreateCollection(db))
	router.GET("/collections/:id/edit", collect.ShowEditCollectionPage(db))
	router.POST("/collections/:id/edit", collect.EditCollection(db))
	router.POST("/collections/:id/delete", collect.DeleteCollection(db))
	log.Printf("Acessing SQLite: %s", config.Database)
	if config.IsTLS {
		log.Printf("Server is running on https://%s", config.TLSListen)
		if err := router.RunTLS(config.TLSListen, config.Cert, config.Key); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	} else {
		router.UseH2C = true
		log.Printf("Server is running on http://%s", config.Listen)
		if err := router.Run(config.Listen); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}
}
