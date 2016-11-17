package conf

import (
	"database/sql"
	"log"
)

var (
	db           *sql.DB
	db_client    *sql.DB
	err          error
	nowWorldTime int64
)

func init() {
	db_url := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Ошибка соединения сервера с БД")
	}
	db_client, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Ошибка соединения клиента с БД")
	}
}
