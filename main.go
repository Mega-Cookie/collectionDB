package main

import (
	"bufio"
	"collectionDB/collect"
	"collectionDB/entries"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zcalusic/sysinfo"
)

var db *sql.DB
var VERSION string
var HOSTNAME string
var info sysinfo.SysInfo
var OS string
var ARCH string
var TZONE string

func getSystemInfo() {
	file, err := os.Open("VERSION")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		VERSION = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	info.GetSysInfo()
	HOSTNAME = info.Node.Hostname
	OS = info.OS.Name
	ARCH = info.Kernel.Architecture
	TZONE = info.Node.Timezone
	_, err = db.Exec(`INSERT OR IGNORE INTO info (HOSTNAME) VALUES (?)`, HOSTNAME)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`UPDATE info SET VERSION = ?, OS = ?, ARCH = ?, TIMEZONE = ?`, VERSION, OS, ARCH, TZONE)
	if err != nil {
		log.Fatal(err)
	}
}
func getStockData() {
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
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./collections.db")
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
		GenreID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME STRING UNIQUE
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS categories (
		CategoryID INTEGER PRIMARY KEY AUTOINCREMENT,
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
		GenreID INT DEFAULT NULL,
		typeID INTEGER DEFAULT NULL,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(collectionID) REFERENCES collections(collectionID),
		FOREIGN KEY(GenreID) REFERENCES genres(GenreID),
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
		CategoryID INT DEFAULT NULL,
		CREATED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		EDITED_AT DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(CategoryID) REFERENCES categories(CategoryID)
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS info (
		instanceID INTEGER PRIMARY KEY,
		VERSION STRING,
		HOSTNAME TEXT UNIQUE,
		OS TEXT,
		ARCH TEXT,
		TIMEZONE TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
func List(db *sql.DB) gin.HandlerFunc {
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
	initDB()
	getSystemInfo()
	getStockData()
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.LoadHTMLGlob("templates/*")
	router.GET("/", List(db))
	router.GET("/create_entry", entries.ShowCreateEntryPage(db))
	router.POST("/create_entry", entries.CreateEntry(db))
	router.GET("/entries/:id/edit", entries.ShowEditEntryPage(db))
	router.POST("/entries/:id/edit", entries.EditEntry(db))
	router.POST("/entries/:id/delete", entries.DeleteEntry(db))
	router.GET("/entries/:id", entries.PreviewSharedEntry(db))
	router.GET("/create_collection", collect.ShowCreateCollectionPage(db))
	router.POST("/create_collection", collect.CreateCollection(db))
	router.GET("/collections/:id/edit", collect.ShowEditCollectionPage(db))
	router.POST("/collections/:id/edit", collect.EditCollection(db))
	router.POST("/collections/:id/delete", collect.DeleteCollection(db))
	port := ":8080"
	log.Printf("Server is running on http://localhost%s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
