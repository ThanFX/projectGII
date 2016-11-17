package client

import (
	"log"
	"strconv"
)

func getTime() {
	var world_time_str string
	var world_time int64
	for {
		err := db.QueryRow("SELECT world_time FROM time WHERE id = 1;").Scan(&world_time_str)
		if err != nil {
			log.Fatal("Ошибка запроса таймеров в БД", err)
		}
		world_time, err = strconv.ParseInt(world_time_str, 10, 64)
	}
}
