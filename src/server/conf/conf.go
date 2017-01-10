package conf

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	ADDR                             string  = ":8080"
	MAX_SOMNOLENCY_FOR_STOP_WORK     float32 = 30.0
	MAX_SOMNOLENCY_FOR_SLEEP         float32 = 40.0
	MIN_SOMNOLENCY_FOR_WAKEUP        float32 = 6.0
	MAX_FATIGUE_FOR_STOP_WORK        float32 = 60.0
	MIN_FATIGUE_FOR_START_WORK       float32 = 20.0
	MAX_MORNING_INTERVAL_BEFORE_WORK int64   = 7200
	MIN_STEP_DURATING                int64   = 600
	MAX_STEP_DURATING                int64   = 3600
)

type Characteristics struct {
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
	PersonId   int         `json: "personId"`
	Name       string      `json: "name"`
	Chunk      PersonChunk `json: "chunk"`
	State      string      `json: "state"`
	Health     float32     `json: "health"`
	Fatigue    float32     `json: "fatigue"`
	Hunger     float32     `json: "hunger"`
	Thirst     float32     `json: "thirst"`
	Somnolency float32     `json: "somnolency"`
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
