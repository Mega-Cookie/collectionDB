package main

import (
	"bufio"
	"collectionDB/entries"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zcalusic/sysinfo"
)

var db *sql.DB
var VERSION string
var hostname string
var info sysinfo.SysInfo
var OS string
var arch string
var tzone string

func getSystemInfo() {
	file, err := os.Open("VERSION")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		VERSION = scanner.Text() // Get the line as a string
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	info.GetSysInfo()

	hostname = info.Node.Hostname
	OS = info.OS.Name
	arch = info.Kernel.Architecture
	tzone = info.Node.Timezone
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./entries.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`SELECT datetime(CURRENT_TIMESTAMP, 'localtime')`)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery := `CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER NOT NULL,
		content TEXT NOT NULL,
		media TEXT NOT NULL,
		is_digital BOOL NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		edited_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	createTableQuery = `CREATE TABLE IF NOT EXISTS info (
		ID INTEGER PRIMARY KEY,
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
	_, err = db.Exec(`INSERT OR IGNORE INTO info (HOSTNAME) VALUES (?)`, hostname)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`UPDATE info SET VERSION = ?, OS = ?, ARCH = ?, TIMEZONE = ?`, VERSION, OS, arch, tzone)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	getSystemInfo()
	initDB()
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.LoadHTMLGlob("templates/*")
	router.GET("/", entries.ListEntries(db))
	router.GET("/create", entries.ShowCreateEntryPage)
	router.POST("/create", entries.CreateEntry(db))
	router.POST("/entries/:id/delete", entries.DeleteEntry(db))
	router.GET("/entries/:id", entries.PreviewSharedEntry(db))
	router.GET("/entries/:id/edit", entries.ShowEditEntryPage(db))
	router.POST("/entries/:id/edit", entries.EditEntry(db))
	port := ":8080"
	log.Printf("Server is running on http://localhost%s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
