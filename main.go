package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"majiang/fourfriends"
	"majiang/proxy"
	"majiang/quedou"
	"sync"
	"time"
)

var TOKEN string

func init() {
	TOKEN = "Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ3eGFwcGxldHRqb285MHI4dHJlYSIsImNyZWF0ZWQiOjE3MTgxNjk1MDc0ODAsImV4cCI6MTgzODE2OTUwN30.pA7sIS8fCAr4UXaCaSQ9KkxeXPFZY848MvAGnyfDpFKLi4I6PGWqWg2CZUzOJay2TKS7Kd2lV2R1sBOSzo0BFw"
}

func initDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"192.168.1.144", "postgres", "pgpassword", "majiang")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db
}

func formatTime() map[string]string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()

	return map[string]string{
		"year":  fmt.Sprintf("%d", year),
		"month": fmt.Sprintf("%d", month),
		"day":   fmt.Sprintf("%d", day),
		"hour":  fmt.Sprintf("%d", hour),
	}
}

func main() {
	proxy.GetProxy()
	panic("test")
	db := initDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	datetime := formatTime()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		quedou.StartWorkOnQuedou(db, datetime, TOKEN)
		wg.Done()
	}()
	go func() {
		fourfriends.StartWorkFourFriends(db, datetime)
		wg.Done()
	}()
	wg.Wait()
}
