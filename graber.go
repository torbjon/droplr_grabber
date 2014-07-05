package main

import (
    "fmt"
    "net/http"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "sync"
    )

var DB *sql.DB
var throttle = make(chan string, 4)

func get_response(url string, wg *sync.WaitGroup, throttle chan string) {
    defer wg.Done()
    uri := strings.Join([]string{"http://d.pr/i/", url}, "")

    transport := http.Transport{ MaxIdleConnsPerHost: 50 }
    client    := http.Client{ Transport: &transport }

    resp, err := client.Get(uri)
    if err != nil {
        // fmt.Println(url, "error")
        log.Print(err)
    } else {
        defer resp.Body.Close()
        err = DB.Ping()
        if err != nil {
            log.Fatalf("Error on opening database connection: %s", err.Error())
        }
        
        var rezult string
        err := DB.QueryRow("SELECT count(id) FROM links WHERE uri=? AND code=?", uri, resp.Status).Scan(&rezult)
        if (err != nil) {
            log.Fatal(err)
        } else if (rezult == "0") {
            _, err = DB.Exec("INSERT INTO links (uri, code) VALUES (?, ?)", url, resp.Status)
            println(url)
        }
    }
    <-throttle
}

func main() {
    chars := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

    db, err := sql.Open("mysql", "root:12345@/graber")
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    DB = db
    DB.SetMaxIdleConns(10)

    var wg sync.WaitGroup

    for _, first_letter := range chars {
        for _, second_letter := range chars {
            for _, third_letter := range chars {
                for _, fourth_letter := range chars {
                  str := strings.Join([]string{first_letter, second_letter, third_letter, fourth_letter}, "")
                  throttle <- str
                  wg.Add(1)
                  go get_response(str, &wg, throttle)
                }
            }
        }
    }
    wg.Wait()

    // get_response("AIF2")
    defer DB.Close()
}