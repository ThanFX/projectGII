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

type CalendarPeriod struct {
	MaxValue      int    `json: "maxValue"`
	MinValue      int    `json: "minValue"`
	PeriodName    string `json: "periodName"`
	PeriodLabel   string `json: "periodLabel"`
	TimeInSeconds int    `json: "timeInSeconds"`
}

var (
	db  *sql.DB
	err error
)

func getCalendarTime(worldTime int) {
	var calendar []CalendarPeriod
	var res string

	//test := `{"maxValue": "60", "minValue": "0", "periodName": "минута", "periodLabel": "minute", "timeInSeconds": "60"}`

	err := db.QueryRow("SELECT value->'periods' FROM config WHERE id = 'calendar';").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса в БД")
	}
	bytes := []byte(res)
	err = json.Unmarshal(bytes, &calendar)
	if err != nil {
		log.Fatal("Ошибка парсинга структуры календаря")
	}

	fmt.Println(calendar)
}

/*
Количество секунд мира на момент старта новой разработки
UPDATE time SET world_time = 439555701 WHERE id = 1;
*/

func set_world_time() {
	var real_time_str, world_time_str, time_speed_str string
	var real_time, world_time int64

	err := db.QueryRow("SELECT real_time, world_time, time_speed FROM time WHERE id = 1").Scan(&real_time_str, &world_time_str, &time_speed_str)
	if err != nil {
		log.Fatal("Ошибка запроса в БД")
	}
	real_time, err = strconv.ParseInt(real_time_str, 10, 64)
	world_time, err = strconv.ParseInt(world_time_str, 10, 64)
	time_speed, err := strconv.Atoi(time_speed_str)

	go create_check(time_speed)
	getCalendarTime(1)

	first_delta := time.Now().Unix() - real_time
	for {
		delta_time := (time.Now().Unix() - real_time) - first_delta
		world_time += (delta_time * int64(time_speed))
		real_time = time.Now().Unix()
		_, err := db.Exec("UPDATE time SET real_time = $1, world_time = $2 WHERE id = 1;", real_time, world_time)
		if err != nil {
			log.Fatal("Ошибка записи в БД")
		}
		time.Sleep(time.Second)
		first_delta = 0
		fmt.Println(real_time, world_time, time_speed)
	}
}

func create_check(world_time_speed int) {
	var res string
	var check CheckPeriods

	err := db.QueryRow("SELECT value::json FROM config WHERE id = 'check_periods'").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса в БД")
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
	} else {
		//defer db.Close()
	}
}

func main() {
	set_world_time()
}
