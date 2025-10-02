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

var VERSION string
var HOSTNAME string
var info sysinfo.SysInfo
var OS string
var ARCH string
var TZONE string
var GOVERSION string

type config struct {
	Listen    string `json:"listen"`
	Database  string `json:"database"`
	Templates string `json:"templates"`
	IsDebug   bool   `json:"debug"`
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
	GOVERSION = runtime.Version()
	_, err = db.Exec(`INSERT OR IGNORE INTO info (HOSTNAME) VALUES (?)`, HOSTNAME)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`UPDATE info SET VERSION = ?, OS = ?, ARCH = ?, GOVERSION = ?, TIMEZONE = ?`, VERSION, OS, ARCH, GOVERSION, TZONE)
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
		fmt.Println(err)
		return
	}
	location, err = time.LoadLocation(zone)
	if err != nil {
		fmt.Println("error: Failed to set timezone")
		fmt.Println(err)
		return
	}
	newstamp = stamp.In(location)
	return
}
