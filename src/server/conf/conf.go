package conf

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	ADDR string = ":8080"
)

var (
	Db        *sql.DB
	Db_client *sql.DB
	err       error
)

func init() {
	db_url := "user=postgres password=postgres dbname=postgres sslmode=disable"
	Db, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Ошибка соединения сервера с БД")
	}
	Db_client, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Ошибка соединения клиента с БД")
	}
}
