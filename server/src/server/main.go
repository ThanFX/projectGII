package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func get_time(db *sql.DB) {
	for {
		cur_time := time.Now().Unix()
		res, err := db.Exec("UPDATE time SET real_time = $1 WHERE id = 1", cur_time)
		fmt.Println(res, err)
		time.Sleep(time.Second)
	}
}

func main() {
	db_url := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Println("Ошибка соединения с БД")
	} else {
		defer db.Close()
	}
	get_time(db)
}
