package small

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/zcalusic/sysinfo"
)

var version string
var hostname string
var info sysinfo.SysInfo
var operatingsystem string
var arch string
var tzone string
var goversion string

type config struct {
	Listen         string `json:"listen"`
	Database       string `json:"database"`
	Templates      string `json:"templates"`
	IsDebug        bool   `json:"debug"`
	IsReverseProxy bool   `json:"proxy"`
	IsTLS          bool   `json:"tls"`
	Cert           string `json:"cert"`
	Key            string `json:"key"`
}
type systeminfo struct {
	ID            int    `json:"id"`
	Version       string `json:"version"`
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	GoVersion     string `json:"go_version"`
	SQLiteVersion string `json:"sqliteversion"`
	Timezone      string `json:"tzone"`
}

func Configure() (config config) {
	filename := flag.String("config", "config.yml", "")
	flag.Parse()
	data, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln(err)
		return
	}
	jconf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf(": Reading config file: %s\n%s", *filename, string(jconf))
	return
}
func SetSystemInfo(db *sql.DB) {
	file, err := os.Open("/etc/collectionDB/VERSION")
	if err != nil {
		fmt.Println("No VERSION file in /etc/collectionDB/ assuming VERSION file in root directory")
		file, err = os.Open("VERSION")
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		version = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	info.GetSysInfo()
	hostname = info.Node.Hostname
	operatingsystem = info.OS.Name
	arch = info.Kernel.Architecture
	tzone = info.Node.Timezone
	goversion = runtime.Version()
	_, err = db.Exec(`INSERT OR IGNORE INTO info (HOSTNAME) VALUES (?)`, hostname)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`UPDATE info SET VERSION = ?, OS = ?, ARCH = ?, TIMEZONE = ?, GOVERSION = ?`, version, operatingsystem, arch, tzone, goversion)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`UPDATE info SET SQLITEVERSION = (SELECT sqlite_version())`, version, operatingsystem, arch, tzone, goversion)
	if err != nil {
		log.Fatal(err)
	}
}
func SetTime(db *sql.DB, stamp *time.Time) (newstamp time.Time) {
	var location *time.Location
	var zone string
	query := "SELECT TIMEZONE FROM info WHERE instanceID = 1"
	err := db.QueryRow(query).Scan(&zone)
	if err != nil {
		fmt.Println("error: Failed to get timezone")
		return
	}
	location, err = time.LoadLocation(zone)
	if err != nil {
		fmt.Println("error: Failed to set timezone")
		return
	}
	newstamp = stamp.In(location)
	return
}
func GetSystemInfo(db *sql.DB) (systeminfo systeminfo) {
	query := "SELECT * FROM info WHERE instanceID = 1"
	response := db.QueryRow(query)
	response.Scan(&systeminfo.ID, &systeminfo.Version, &systeminfo.Hostname, &systeminfo.OS, &systeminfo.Arch, &systeminfo.GoVersion, &systeminfo.SQLiteVersion, &systeminfo.Timezone)
	return
}
