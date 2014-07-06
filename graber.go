package main

import (
	"database/sql"
	_ "fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	DB       *sql.DB
	throttle = make(chan string, 4)
	site_url = "http://d.pr/i/"
	chars    = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
)

func main() {
	db_connection()

	var wg sync.WaitGroup

	for _, first_letter := range chars {
		for _, second_letter := range chars {
			for _, third_letter := range chars {
				for _, fourth_letter := range chars {
					str := strings.Join([]string{first_letter, second_letter, third_letter, fourth_letter}, "")
					throttle <- str
					wg.Add(1)
					go grab_url(str, &wg, throttle)
				}
			}
		}
	}
	// grab_url("AAAS", &wg, throttle)
	wg.Wait()

	defer DB.Close()
}

func grab_url(url string, wg *sync.WaitGroup, throttle chan string) {
	defer wg.Done()

	if db_check_uri_code(url, uint64(200)) == "0" {
		uri := strings.Join([]string{site_url, url}, "")
		transport := http.Transport{MaxIdleConnsPerHost: 50}
		client := http.Client{Transport: &transport}
		resp, err := client.Get(uri)
		if err != nil {
			log.Print(err)
		} else {
			defer resp.Body.Close()
			db_ping()
			db_insert_link(url, resp.Status)
			println(url)
		}
	}
	<-throttle
}

func db_connection() {
	/* "CREATE TABLE `links` (
	   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	   `uri` varchar(10) DEFAULT '',
	   `code` smallint(1) DEFAULT NULL,
	   PRIMARY KEY (`id`)
	 ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;" */

	db, err := sql.Open("mysql", "root:12345@/graber")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	DB = db
	DB.SetMaxIdleConns(10)
}

func db_ping() {
	var err = DB.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
}

func db_check_uri(url string) string {
	var rezult string
	err := DB.QueryRow("SELECT count(id) FROM links WHERE uri = ?", url).Scan(&rezult)
	if err != nil {
		log.Fatal(err)
	}
	return rezult
}

func db_check_uri_code(url string, code uint64) string {
	var rezult string
	err := DB.QueryRow("SELECT count(id) FROM links WHERE uri = ? AND code = ?", url, code).Scan(&rezult)
	if err != nil {
		log.Fatal(err)
	}
	return rezult
}

func db_insert_link(url string, status string) {
	var code = response_status_to_code(status)
	if db_check_uri(url) == "0" {
		_, _ = DB.Exec("INSERT INTO links (uri, code) VALUES (?, ?)", url, code)
	} else {
		_, _ = DB.Exec("UPDATE links SET code = ? WHERE uri = ?", code, url)
	}
}

func response_status_to_code(status string) uint64 {
	var code = strings.Split(status, " ")
	d, _ := strconv.ParseUint(code[0], 0, 16)
	return d
}
