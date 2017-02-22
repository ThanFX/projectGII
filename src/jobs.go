package main

import (
	"encoding/json"
	"log"
	"math"
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
	Work  int
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

	err := db.QueryRow(`SELECT value::json FROM config WHERE id = 'check_periods'`).Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса периодов проверки в БД")
	}

	init_check()
	prepareQueries()

	bytes := []byte(res)
	json.Unmarshal(bytes, &check)

	// TODO Добавить проверку на минимальный шаг запуска в 1 секунду!!
	state_period := strconv.Itoa(int(math.Max(float64(check.State*60/world_time_speed), 1)))
	state_period = "@every " + state_period + "s"
	hts_period := strconv.Itoa(int(math.Max(float64(check.HTS*60/world_time_speed), 1)))
	hts_period = "@every " + hts_period + "s"
	ed_period := strconv.Itoa(int(math.Max(float64(check.ED*60/world_time_speed), 1)))
	ed_period = "@every " + ed_period + "s"
	work_period := strconv.Itoa(int(math.Max(float64(check.Work*60/world_time_speed), 1)))
	work_period = "@every " + work_period + "s"

	state_cron := cron.New()
	state_cron.AddFunc(state_period, state_job)
	state_job()
	state_cron.Start()

	hts_cron := cron.New()
	hts_cron.AddFunc(hts_period, hts_job)
	hts_job()
	hts_cron.Start()

	ed_cron := cron.New()
	ed_cron.AddFunc(ed_period, ed_job)
	ed_job()
	ed_cron.Start()

	task_cron := cron.New()
	task_cron.AddFunc(work_period, step_job)
	step_job()
	task_cron.Start()
}

func state_job() {
	if lib.GetNowWorldTime() == 0 {
		return
	}

	cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	//fmt.Println("Проверка состояний началась в ", lib.GetWCTString(cTime))

	cH, err := strconv.Atoi(cTime["hour"])
	if err != nil {
		log.Fatal("Ошибка парсинга времени: ", err)
	}

	if cH < 6 || cH >= 20 {
		_, err := db.Exec(`UPDATE persons SET state = $1 WHERE state = $2 AND somnolency > $3;`,
			"sleep", "chores", conf.MAX_SOMNOLENCY_FOR_SLEEP)
		if err != nil {
			log.Fatal("Ошибка засыпания персонажей: ", err)
		}
	}
	if cH >= 6 && cH < 20 {
		upd, err := db.Exec(`UPDATE persons SET state = $1 WHERE state = $2 AND somnolency <= $3;`,
			"chores", "sleep", conf.MIN_SOMNOLENCY_FOR_WAKEUP)
		if err != nil {
			log.Fatal("Ошибка пробуждения персонажей: ", err)
		}
		countPerson, err := upd.RowsAffected()
		if err != nil {
			log.Fatal("Ошибка получения количества проснувшихся персонажей")
		}
		// Если хоть кто-то проснулся, есть смысл запустить создание работ для них
		if countPerson > 0 {
			log.Printf("Проснулось %d персонажей\n", countPerson)
			create_task()
		}
	}

}

func ed_job() {
	if lib.GetNowWorldTime() == 0 {
		return
	}

	//cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	//fmt.Println("Едим и пьём в ", lib.GetWCTString(cTime))

	_, err := db.Exec(`UPDATE persons SET thirst = $1 WHERE thirst >= 6.0 AND state != $2;`, 0.0, "sleep")
	if err != nil {
		log.Fatal("Ошибка утоления жажды: ", err)
	}

	_, err = db.Exec(`
		UPDATE persons SET thirst = $1, hunger = $2 WHERE
			hunger >= 3.0 AND state in ($3, $4);
		`, 0.0, 0.0, "chores", "rest")
	if err != nil {
		log.Fatal("Ошибка утоления голода: ", err)
	}
}

func hts_job() {
	if lib.GetNowWorldTime() == 0 {
		return
	}
	_, err := queries["updateAllPersonsHTFS"].query.Exec(
		stateSpeeds.Get("sleep.hunger").Float64(),
		stateSpeeds.Get("sleep.thirst").Float64(),
		stateSpeeds.Get("sleep.fatigue").Float64(),
		stateSpeeds.Get("sleep.somnolency").Float64(),
		lib.GetNowWorldTime(), "sleep")
	if err != nil {
		log.Fatalf(queries["updateAllPersonsHTFS"].queryErrorText, "sleep", err)
	}
	_, err = queries["updateAllPersonsHTFS"].query.Exec(
		stateSpeeds.Get("chores.hunger").Float64(),
		stateSpeeds.Get("chores.thirst").Float64(),
		stateSpeeds.Get("chores.fatigue").Float64(),
		stateSpeeds.Get("chores.somnolency").Float64(),
		lib.GetNowWorldTime(), "chores")
	if err != nil {
		log.Fatalf(queries["updateAllPersonsHTFS"].queryErrorText, "chores", err)
	}
	_, err = queries["updateAllPersonsHTFS"].query.Exec(
		stateSpeeds.Get("work.hunger").Float64(),
		stateSpeeds.Get("work.thirst").Float64(),
		stateSpeeds.Get("work.fatigue").Float64(),
		stateSpeeds.Get("work.somnolency").Float64(),
		lib.GetNowWorldTime(), "work")
	if err != nil {
		log.Fatalf(queries["updateAllPersonsHTFS"].queryErrorText, "work", err)
	}
	_, err = queries["updateAllPersonsHTFS"].query.Exec(
		stateSpeeds.Get("rest.hunger").Float64(),
		stateSpeeds.Get("rest.thirst").Float64(),
		stateSpeeds.Get("rest.fatigue").Float64(),
		stateSpeeds.Get("rest.somnolency").Float64(),
		lib.GetNowWorldTime(), "rest")
	if err != nil {
		log.Fatalf(queries["updateAllPersonsHTFS"].queryErrorText, "rest", err)
	}
}
