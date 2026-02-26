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
		DESCRIPTION TEXT,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TABLE IF NOT EXISTS casetypes (
		casetypeID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		DESCRIPTION TEXT,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TABLE IF NOT EXISTS genres (
		genreID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		DESCRIPTION TEXT,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TABLE IF NOT EXISTS categories (
		categoryID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		DESCRIPTION TEXT,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TABLE IF NOT EXISTS publishers (
		publisherID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE,
		DESCRIPTION TEXT,
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
		categoryID INTEGER NOT NULL DEFAULT 1,
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

	createTableQuery = `CREATE TABLE IF NOT EXISTS entries (
		entryID INTEGER PRIMARY KEY AUTOINCREMENT,
		TITLE TEXT NOT NULL,
		YEAR INTEGER NOT NULL,
		PLOT TEXT NOT NULL,
		COMMENT TEXT NOT NULL,
		AUDIOLANGS TEXT,
		SUBTITLELANGS TEXT,
		MEDIARELEASEDATE TEXT,
		MEDIACOUNT INTEGER,
		IS_DIGITAL BOOLEAN,
		IS_BOOKLET BOOLEAN,
		REGIONCODE INTEGER,
		BARCODE TEXT,
		collectionID INTEGER NOT NULL,
		genreID INTEGER NOT NULL,
		mediatypeID INTEGER NOT NULL,
		casetypeID INTEGER NOT NULL,
		publisherID INTEGER NOT NULL,
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

	createTableQuery = `CREATE TRIGGER IF NOT EXISTS prevent_delete_default_collection
		BEFORE DELETE ON collections
		FOR EACH ROW
		WHEN OLD.collectionID = 1
		BEGIN
    	SELECT RAISE(ABORT, 'You cant delete the default Collection.');
		END;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TRIGGER IF NOT EXISTS prevent_delete_default_category
		BEFORE DELETE ON categories
		FOR EACH ROW
		WHEN OLD.categoryID = 1
		BEGIN
    	SELECT RAISE(ABORT, 'You cant delete the default Category.');
		END;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TRIGGER IF NOT EXISTS prevent_delete_default_genre
		BEFORE DELETE ON genres
		FOR EACH ROW
		WHEN OLD.genreID = 1
		BEGIN
    	SELECT RAISE(ABORT, 'You cant delete the default Genre.');
		END;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TRIGGER IF NOT EXISTS prevent_delete_default_case
		BEFORE DELETE ON casetypes
		FOR EACH ROW
		WHEN OLD.casetypeID = 1
		BEGIN
    	SELECT RAISE(ABORT, 'You cant delete the default Case Type.');
		END;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TRIGGER IF NOT EXISTS prevent_delete_default_media
		BEFORE DELETE ON mediatypes
		FOR EACH ROW
		WHEN OLD.mediatypeID = 1
		BEGIN
    	SELECT RAISE(ABORT, 'You cant delete the default Media Type.');
		END;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery = `CREATE TRIGGER IF NOT EXISTS prevent_delete_default_publisher
		BEFORE DELETE ON publishers
		FOR EACH ROW
		WHEN OLD.publisherID = 1
		BEGIN
    	SELECT RAISE(ABORT, 'You cant delete the default Publisher.');
		END;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
func GetCollections(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Collections",
			"data": gin.H{
				"Collections": collect.ListCollections(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetEntries(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Entries",
			"data": gin.H{
				"Entries": entries.ListEntries(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetCaseTypes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Case Types",
			"data": gin.H{
				"CaseTypes": stockdata.ListCaseTypes(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetMediaTypes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Media Types",
			"data": gin.H{
				"MediaTypes": stockdata.ListMediaTypes(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Categories",
			"data": gin.H{
				"Categories": stockdata.ListCategories(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetGenres(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Genres",
			"data": gin.H{
				"Genres": stockdata.ListGenres(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetPublishers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Categories",
			"data": gin.H{
				"Publishers": stockdata.ListPublishers(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func GetAbout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		answer := gin.H{
			"Status":  http.StatusOK,
			"Message": "Successfully loaded Systeminfo",
			"data": gin.H{
				"Info": small.GetSystemInfo(db)},
		}
		c.JSON(http.StatusOK, answer)
	}
}
func ShowStock() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "stock/index.html", gin.H{})
	}
}
func ShowIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}
func ShowAbout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", gin.H{})
	}
}
func main() {
	config := small.Configure()
	initDB(config.Database)
	small.SetStockData(db)
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
	// Browser
	router.GET("/", ShowIndex())
	router.GET("/stock", ShowStock())
	router.GET("/about", ShowAbout())
	router.POST("/stock/mediatype/create", stockdata.CreateMediaType(db))
	router.POST("/stock/casetype/create", stockdata.CreateCaseType(db))
	router.POST("/stock/publisher/create", stockdata.CreatePublisher(db))
	router.POST("/stock/category/create", stockdata.CreateCategory(db))
	router.POST("/stock/genre/create", stockdata.CreateGenre(db))
	router.DELETE("/publisher/:id/delete", stockdata.DeletePublisher(db))
	router.DELETE("/category/:id/delete", stockdata.DeleteCategory(db))
	router.DELETE("/genre/:id/delete", stockdata.DeleteGenre(db))
	router.GET("/entries/:id", entries.ViewEntry(db))
	router.GET("/create_entry", entries.ShowCreateEntryPage(db))
	router.POST("/create_entry", entries.CreateEntry(db))
	router.GET("/entries/:id/edit", entries.ShowEditEntryPage(db))
	router.POST("/entries/:id/edit", entries.EditEntry(db))
	router.GET("/collections/:id", collect.ViewCollection(db))
	router.GET("/create_collection", collect.ShowCreateCollectionPage(db))
	router.POST("/create_collection", collect.CreateCollection(db))
	router.GET("/collections/:id/edit", collect.ShowEditCollectionPage(db))
	router.POST("/collections/:id/edit", collect.EditCollection(db))
	// API
	router.GET("/api/v1/entries", GetEntries(db))
	router.GET("/api/v1/collections", GetCollections(db))
	router.GET("/api/v1/casetypes", GetCaseTypes(db))
	router.GET("/api/v1/mediatypes", GetMediaTypes(db))
	router.GET("/api/v1/categories", GetCategories(db))
	router.GET("/api/v1/genres", GetGenres(db))
	router.GET("/api/v1/publishers", GetPublishers(db))
	router.DELETE("/api/v1/collection/:id", collect.DeleteCollection(db))
	router.DELETE("/api/v1/entry/:id", entries.DeleteEntry(db))
	router.DELETE("/api/v1/mediatype/:id", stockdata.DeleteMediaType(db))
	router.DELETE("/api/v1/casetype/:id", stockdata.DeleteCaseType(db))
	router.GET("/api/v1/about", GetAbout(db))
	// log
	log.Printf("Accessing SQLite: %s", config.Database)
	// SSL
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
