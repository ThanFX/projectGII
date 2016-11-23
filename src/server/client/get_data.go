package client

import (
	"log"
	"server/lib"
	"strconv"
	"time"
)

func getTime(sendTime chan []byte) {
	var world_time_str string
	var world_time int64
	for {
		err := db.QueryRow("SELECT world_time FROM time WHERE id = 1;").Scan(&world_time_str)
		if err != nil {
			log.Fatal("Ошибка запроса таймеров в БД", err)
		}
		world_time, err = strconv.ParseInt(world_time_str, 10, 64)
		//fmt.Println("Cчитанное из БД время: " + lib.GetWCTString(lib.GetWorldCalendarTime(world_time)))
		sendTime <- []byte(lib.GetWCTString(lib.GetWorldCalendarTime(world_time)))
		time.Sleep(time.Second)
	}
}
