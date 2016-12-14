package conf

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	ADDR string = ":8080"
)

type Characteristics struct {
	State      int     `json: "state"`
	Health     float32 `json: "health"`
	Fatigue    float32 `json: "fatigue"`
	Hunger     float32 `json: "hunger"`
	Thirst     float32 `json: "thirst"`
	Somnolency float32 `json: "somnolency"`
}

type PersonChunk struct {
	X int `json: "x"`
	Y int `json: "y"`
}

type Person struct {
	PersonId  int             `json: "personId"`
	Name      string          `json: "name"`
	Chunk     PersonChunk     `json: "chunk"`
	PersonChr Characteristics `json: "characteristics"`
}

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
