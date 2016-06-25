package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/robfig/cron"
)

type StateSpeed struct {
	State string          `json: "state"`
	Speed Characteristics `json: "speed"`
}

type Characteristics struct {
	Fatigue    float32
	Hunger     float32
	Thirst     float32
	Somnolency float32
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
	stateSpeeds []StateSpeed
)

func init_check() {
	var res string
	err := db.QueryRow("SELECT value->'states' FROM config WHERE id = 'state_speed';").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса скоростей изменений характеристик в БД", err)
	}
	bytes := []byte(res)
	err = json.Unmarshal(bytes, &stateSpeeds)
	if err != nil {
		log.Fatal("Ошибка парсинга структуры скоростей изменений характеристик", err)
	}
}

func create_check(world_time_speed int) {
	var res string
	var check CheckPeriods

	err := db.QueryRow("SELECT value::json FROM config WHERE id = 'check_periods'").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса периодов проверки в БД")
	}

	init_check()

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
	cTime := getWorldCalendarTime(nowWorldTime)
	fmt.Println("Проверка состояний началась в ", getWCTString(cTime))

}

func ed_job() {
	cTime := getWorldCalendarTime(nowWorldTime)
	fmt.Println("Едим и пьём в ", getWCTString(cTime))
}

func hts_job() {
	cTime := getWorldCalendarTime(nowWorldTime)
	fmt.Println("Проверка работы началась в ", getWCTString(cTime))

	_, err := db.Exec(`
		UPDATE person_health_characteristic SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		;`,
		5,
		10,
		-10,
		-20,
		nowWorldTime)
	if err != nil {
		log.Fatal("Ошибка обновления характерстик: ", err)
	}
}
