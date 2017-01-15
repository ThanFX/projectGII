package main

import (
	"log"
	"server/conf"
	"server/lib"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var (
	db = conf.Db
)

/*
Количество секунд мира на момент старта новой разработки
UPDATE time SET world_time = 439555701 WHERE id = 1;
*/

func setWorldTime() {
	var real_time_str, world_time_str, time_speed_str string
	var real_time, world_time int64

	err := db.QueryRow("SELECT real_time, world_time, time_speed FROM time WHERE id = 1;").Scan(&real_time_str, &world_time_str, &time_speed_str)
	if err != nil {
		log.Fatal("Ошибка запроса таймеров в БД", err)
	}
	real_time, err = strconv.ParseInt(real_time_str, 10, 64)
	world_time, err = strconv.ParseInt(world_time_str, 10, 64)
	time_speed, err := strconv.Atoi(time_speed_str)

	go create_check(time_speed)

	first_delta := time.Now().Unix() - real_time
	for {
		delta_time := (time.Now().Unix() - real_time) - first_delta
		world_time += (delta_time * int64(time_speed))
		real_time = time.Now().Unix()
		_, err := db.Exec("UPDATE time SET real_time = $1, world_time = $2 WHERE id = 1;", real_time, world_time)
		if err != nil {
			log.Fatal("Ошибка записи таймеров в БД")
		}
		lib.SetNowWorldTime(world_time)
		first_delta = 0
		time.Sleep(time.Second)
	}
}

func main() {
	defer db.Close()
	lib.GetCalendar()
	//go client.ClientStart()
	//setWorldTime()
	createWorkResult()
}
