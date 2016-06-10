package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

/*
Количество секунд мира на момент старта новой разработки
UPDATE time SET world_time = 439555701 WHERE id = 1;
*/

func set_world_time(db *sql.DB) {
	var real_time_str, world_time_str, time_speed_str string
	var real_time, world_time int64

	err := db.QueryRow("SELECT real_time, world_time, time_speed FROM time WHERE id = 1").Scan(&real_time_str, &world_time_str, &time_speed_str)
	if err != nil {
		log.Fatal("Ошибка запроса в БД")
	}
	real_time, err = strconv.ParseInt(real_time_str, 10, 64)
	world_time, err = strconv.ParseInt(world_time_str, 10, 64)
	time_speed, err := strconv.Atoi(time_speed_str)
	first_delta := time.Now().Unix() - real_time
	for {
		delta_time := (time.Now().Unix() - real_time) - first_delta
		world_time += (delta_time * int64(time_speed))
		real_time = time.Now().Unix()
		_, err := db.Exec("UPDATE time SET real_time = $1, world_time = $2 WHERE id = 1", real_time, world_time)
		if err != nil {
			log.Fatal("Ошибка записи в БД")
		}
		time.Sleep(time.Second)
		first_delta = 0
		fmt.Println(real_time, world_time, time_speed)
	}
}

func main() {
	db_url := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Ошибка соединения с БД")
	} else {
		defer db.Close()
	}
	set_world_time(db)
}
