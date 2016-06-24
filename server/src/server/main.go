package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"gopkg.in/robfig/cron.v2"
)

type Characteristics struct {
	health     float32
	fatigue    float32
	hunger     float32
	thirst     float32
	somnolency float32
}

type Person struct {
	person_id int
	name      string
	chr       Characteristics
}

type CheckPeriods struct {
	ED    int
	HTS   int
	State int
}

var (
	db  *sql.DB
	err error
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
		//go updateWorldCalendarTime(world_time)
		first_delta = 0
		//fmt.Println(real_time, world_time, time_speed)
		fmt.Println(getWCTString(getWorldCalendarTime(world_time)))
		time.Sleep(time.Second)
	}
}

func create_check(world_time_speed int) {
	var res string
	var check CheckPeriods

	err := db.QueryRow("SELECT value::json FROM config WHERE id = 'check_periods'").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса периодов проверки в БД")
	}
	bytes := []byte(res)
	json.Unmarshal(bytes, &check)

	state_period := strconv.Itoa(int((check.State * 60) / world_time_speed))
	state_period = "@every " + state_period + "s"
	hts_period := strconv.Itoa(int((check.HTS * 60) / world_time_speed))
	hts_period = "@every " + hts_period + "s"
	ed_period := strconv.Itoa(int((check.ED * 60) / world_time_speed))
	ed_period = "@every " + ed_period + "s"

	state_cron := cron.New()
	state_cron.AddFunc(state_period, state_job)
	go state_job()
	state_cron.Start()

	hts_cron := cron.New()
	hts_cron.AddFunc(hts_period, hts_job)
	go hts_job()
	hts_cron.Start()

	ed_cron := cron.New()
	ed_cron.AddFunc(ed_period, ed_job)
	go ed_job()
	ed_cron.Start()
}

func state_job() {
	fmt.Println("Проверка состояний")
}

func ed_job() {
	fmt.Println("Едим и пьём")
}

func hts_job() {
	fmt.Println("Проверка голода, жажды и сонливости")
}

func init() {
	db_url := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Ошибка соединения с БД")
	}
}

func main() {
	defer db.Close()
	getCalendar()
	setWorldTime()
}
