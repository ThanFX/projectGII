package main

import (
	"encoding/json"
	"fmt"
	"log"
	"server/conf"
	"server/lib"
	"strconv"

	"github.com/robfig/cron"
	"github.com/stretchr/objx"
)

type StateSpeed struct {
	State string               `json: "state"`
	Speed conf.Characteristics `json: "speed"`
}

type CheckPeriods struct {
	ED    int
	HTS   int
	State int
}

var (
	stateSpeeds objx.Map
)

func init_check() {
	var res string
	err := db.QueryRow("SELECT value FROM config WHERE id = 'state_speed';").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса скоростей изменений характеристик в БД", err)
	}
	stateSpeeds, err = objx.FromJSON(res)
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
	if lib.GetNowWorldTime() == 0 {
		return
	}

	cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	fmt.Println("Проверка состояний началась в ", lib.GetWCTString(cTime))

	cH, err := strconv.Atoi(cTime["hour"])
	if err != nil {
		log.Fatal("Ошибка парсинга времени: ", err)
	}

	if cH < 6 || cH >= 20 {
		_, err := db.Exec(`
		UPDATE person_health_characteristic SET state = $1 WHERE
			somnolency > 40.0;
		`, 1)
		if err != nil {
			log.Fatal("Ошибка засыпания персонажей: ", err)
		}
	}
	if cH >= 6 && cH < 20 {
		_, err := db.Exec(`
		UPDATE person_health_characteristic SET state = $1 WHERE
			somnolency <= 6.0;
		`, 5)
		if err != nil {
			log.Fatal("Ошибка пробуждения персонажей: ", err)
		}
	}
}

func ed_job() {
	if lib.GetNowWorldTime() == 0 {
		return
	}

	cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	fmt.Println("Едим и пьём в ", lib.GetWCTString(cTime))

	_, err := db.Exec(`
		UPDATE person_health_characteristic SET thirst = $1 WHERE
			thirst >= 6.0 AND state != $2;
		`, 0.0, 1)
	if err != nil {
		log.Fatal("Ошибка утоления жажды: ", err)
	}

	_, err = db.Exec(`
		UPDATE person_health_characteristic SET thirst = $1, hunger = $2 WHERE
			hunger >= 3.0 AND state != $3;
		`, 0.0, 0.0, 1)
	if err != nil {
		log.Fatal("Ошибка утоления голода: ", err)
	}
}

func hts_job() {
	if lib.GetNowWorldTime() == 0 {
		return
	}

	cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	fmt.Println("Проверка работы началась в ", lib.GetWCTString(cTime))

	_, err := db.Exec(`
		UPDATE person_health_characteristic SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = 1
		;`,
		stateSpeeds.Get("sleep.hunger").Float64(),
		stateSpeeds.Get("sleep.thirst").Float64(),
		stateSpeeds.Get("sleep.fatigue").Float64(),
		stateSpeeds.Get("sleep.somnolency").Float64(),
		lib.GetNowWorldTime())
	if err != nil {
		log.Fatal("Ошибка обновления характерстик спящих персонажей: ", err)
	}
	_, err = db.Exec(`
		UPDATE person_health_characteristic SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = 5
		;`,
		stateSpeeds.Get("chores.hunger").Float64(),
		stateSpeeds.Get("chores.thirst").Float64(),
		stateSpeeds.Get("chores.fatigue").Float64(),
		stateSpeeds.Get("chores.somnolency").Float64(),
		lib.GetNowWorldTime())
	if err != nil {
		log.Fatal("Ошибка обновления характерстик бодрствующих персонажей: ", err)
	}
}
