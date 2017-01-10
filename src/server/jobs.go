package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"server/conf"
	"server/lib"
	"strconv"
	"strings"
	"time"

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

type Query struct {
	query          *sql.Stmt
	text           string
	queryErrorText string
}

var (
	stateSpeeds objx.Map
	queries     map[string]Query
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

func initQueries() {
	var emptyStmp *sql.Stmt
	queries = make(map[string]Query)
	// $1 - startDayTime, $2 - nowTime
	queries["getPersonWithoutWork"] = Query{
		emptyStmp,
		`SELECT p.id
  		FROM
    		(SELECT id
    		FROM persons
      		WHERE state = 'chores') p
    	LEFT JOIN
    		(SELECT person_id FROM tasks
    		WHERE create_time > $1 AND create_time < $2) t ON t.person_id = p.id
		WHERE t.person_id IS NULL;`,
		`Ошибка выполнения запроса получения списка персонажей без задач на сегодня: %s`}
	// $1 - person_id
	queries["getPreferSkill"] = Query{
		emptyStmp,
		`SELECT skill_id FROM person_skills WHERE person_id = $1 ORDER BY worth DESC LIMIT 1;`,
		`Ошибка получения наилучшей работы для персонажа %d: %s`}
	//
	queries["createTask"] = Query{
		emptyStmp,
		`INSERT INTO tasks(person_id, skill_id, steps, result, create_time) VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		`Ошибка создания задачи на работу для персонажа %d в %d: $s`}
	// $1 - state, $2 - personId
	queries["setPersonState"] = Query{
		emptyStmp,
		`UPDATE persons SET state = $1 WHERE id = $2;`,
		`Ошибка установки состояния \"%s\" для персонажа %d: %s`}
	//$1 - taskId, $2 - step
	queries["setStepDone"] = Query{
		emptyStmp,
		`UPDATE task_steps SET is_done = TRUE WHERE task_id = $1 AND step = $2;`,
		`Ошибка закрытия шага %d для задачи %d: %s`}
	//$1 - taskId, $2 - step
	queries["getStepData"] = Query{
		emptyStmp,
		`SELECT (steps->$2)->'finish_time' from tasks WHERE id = $1;`,
		`Ошибка получения данных шага %d задачи %d: %s`}
	//$1 - taskId
	queries["setTaskDone"] = Query{
		emptyStmp,
		`UPDATE tasks SET is_done = TRUE WHERE id = $1;`,
		`Ошибка закрытия задачи %d: %s`}
	//
	queries["createNewStep"] = Query{
		emptyStmp,
		`INSERT INTO task_steps(task_id, step, start_time, finish_time, type, create_time) VALUES ($1, $2, $3, $4, $5, $6);`,
		`Ошибка создания шага %d задачи %d: %s`}
	//
	queries["getFinishedSteps"] = Query{
		emptyStmp,
		`SELECT ts.task_id, ts.step, ts.type, t.person_id, ts.finish_time, p.fatigue, p.somnolency
				FROM task_steps ts
				JOIN tasks t on ts.task_id = t.id
				JOIN persons p ON p.id = t.person_id
				WHERE ts.is_done = FALSE AND ts.finish_time < $1;`,
		`Ошибка получения завершившихся шагов: %s`}
}

func prepareQueries() {
	initQueries()
	var q Query
	for key, _ := range queries {
		prepareQuery, err := db.Prepare(queries[key].text)
		if err != nil {
			log.Fatal("Ошибка подготовки запроса ", key, ": ", err)
		}
		q.query = prepareQuery
		q.text = queries[key].text
		q.queryErrorText = queries[key].queryErrorText
		queries[key] = q
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
		_, err := db.Exec(`
		UPDATE persons SET state = $1 WHERE
			somnolency > 40.0;
		`, "sleep")
		if err != nil {
			log.Fatal("Ошибка засыпания персонажей: ", err)
		}
	}
	if cH >= 6 && cH < 20 {
		upd, err := db.Exec(`
		UPDATE persons SET state = $1 WHERE
			somnolency <= 6.0;
		`, "chores")
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

	//cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	//fmt.Println("Проверка работы началась в ", lib.GetWCTString(cTime))

	_, err := db.Exec(`
		UPDATE persons SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = 'sleep'
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
		UPDATE persons SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = 'chores'
		;`,
		stateSpeeds.Get("chores.hunger").Float64(),
		stateSpeeds.Get("chores.thirst").Float64(),
		stateSpeeds.Get("chores.fatigue").Float64(),
		stateSpeeds.Get("chores.somnolency").Float64(),
		lib.GetNowWorldTime())
	if err != nil {
		log.Fatal("Ошибка обновления характерстик персонажей, занимающихся домашними делами: ", err)
	}
	_, err = db.Exec(`
		UPDATE persons SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = 'work'
		;`,
		stateSpeeds.Get("work.hunger").Float64(),
		stateSpeeds.Get("work.thirst").Float64(),
		stateSpeeds.Get("work.fatigue").Float64(),
		stateSpeeds.Get("work.somnolency").Float64(),
		lib.GetNowWorldTime())
	if err != nil {
		log.Fatal("Ошибка обновления характерстик работающих персонажей: ", err)
	}
	_, err = db.Exec(`
		UPDATE persons SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = 'rest'
		;`,
		stateSpeeds.Get("rest.hunger").Float64(),
		stateSpeeds.Get("rest.thirst").Float64(),
		stateSpeeds.Get("rest.fatigue").Float64(),
		stateSpeeds.Get("rest.somnolency").Float64(),
		lib.GetNowWorldTime())
	if err != nil {
		log.Fatal("Ошибка обновления характерстик отдыхающих персонажей: ", err)
	}
}

// Создаём задачи (и массив шагов выполнения для задач) на старт работ для персонажей
func create_task() {
	var personId, taskId int
	nowTime := lib.GetNowWorldTime()
	startDayTime := lib.GetStartDayTime(nowTime)
	// Получаем список персонажей, у которых нет работы
	//fmt.Println("getPersonWithoutWork", startDayTime)
	persons, err := queries["getPersonWithoutWork"].query.Query(startDayTime, nowTime)
	//fmt.Println("!")
	if err != nil {
		log.Fatalf(queries["getPersonWithoutWork"].queryErrorText, err)
	}
	defer persons.Close()
	rand.Seed(time.Now().UTC().UnixNano())

	for persons.Next() {
		err = persons.Scan(&personId)
		if err != nil {
			log.Fatal("Ошибка парсинга списка персонажей: ", err)
		}
		//fmt.Println("Получаем список персонажей - ", personId)
		preferSkillId := getPreferPersonSkill(personId)
		nowTime := lib.GetNowWorldTime()
		steps := getTaskSteps(nowTime)
		//fmt.Println("createTask")
		err = queries["createTask"].query.QueryRow(personId, preferSkillId, steps, "{}", nowTime).Scan(&taskId)
		//fmt.Println("!")
		if err != nil {
			log.Fatalf(queries["createTask"].queryErrorText, personId, startDayTime, err)
		}
		createNewStep("wait", 0, taskId, int(nowTime))
		fmt.Println("Создали работу для ", personId)
	}
}

func getTaskSteps(nowTime int64) string {
	// Первый шаг через некоторый промежуток после пробуждения
	stepTime := rand.Int63n(conf.MAX_MORNING_INTERVAL_BEFORE_WORK) + conf.MIN_STEP_DURATING + nowTime
	steps := "{"
	step := 0
	// Создаём массив шагов, до конца текущих суток, нулевой шаг - ошидание работы с момента пробуждения,
	// после него начинается реальная работа
	for stepTime < (lib.GetStartDayTime(nowTime) + 24*3600) {
		steps += "\"" + strconv.Itoa(step) + "\":{\"finish_time\":" + strconv.FormatInt(stepTime, 10)
		steps += "},"
		step++
		stepTime += rand.Int63n(conf.MAX_STEP_DURATING-conf.MIN_STEP_DURATING) + conf.MIN_STEP_DURATING
	}
	steps = strings.TrimRight(steps, ",")
	steps += "}"
	return steps
}

func getPreferPersonSkill(personId int) int {
	var skillId int
	//fmt.Println("getPreferSkill")
	err := queries["getPreferSkill"].query.QueryRow(personId).Scan(&skillId)
	if err != nil {
		log.Fatalf(queries["getPreferSkill"].queryErrorText, personId, err)
	}
	return skillId
}

func step_job() {
	var taskId, personId, step, finish_time int
	var fatifue, somnolency float32
	var taskType string

	nowTime := lib.GetNowWorldTime()
	//fmt.Println("getFinishedSteps")
	tasksList, err := queries["getFinishedSteps"].query.Query(nowTime)
	//fmt.Println("!")
	if err != nil {
		log.Fatal("Ошибка получения наступивших задач на шаги работы: ", err)
	}
	defer tasksList.Close()

	for tasksList.Next() {
		// ts.task_id, ts.step, ts.type, t.person_id, ts.finish_time, chr.fatigue, chr.somnolency
		err = tasksList.Scan(&taskId, &step, &taskType, &personId, &finish_time, &fatifue, &somnolency)
		if err != nil {
			log.Fatal("Ошибка парсинга списка наступивших задач на шаги работы: ", err)
		}
		//fmt.Println("смотрим персонажа ", personId, " задача ", taskId)
		// Если шаг нулевой - закрываем его и начинаем честную работу
		if step == 0 {
			setStepDone(taskId, step)
			createNewStep("work", step+1, taskId, finish_time)
			setPersonState(personId, "work")
			log.Printf("%d | Персонаж %d начал работать, сонливость - %.2f\n", nowTime, personId, somnolency)
			// Обрабатываем окончание работы - смотрим на повышенную сонливость
		} else if somnolency > conf.MAX_SOMNOLENCY_FOR_STOP_WORK {
			// Если не устали и при этом работали - продолжаем пахать
			if fatifue < conf.MAX_FATIGUE_FOR_STOP_WORK && taskType == "work" {
				setStepDone(taskId, step)
				createNewStep("work", step+1, taskId, finish_time)
			} else {
				// Если же устали или уже отдыхали на работе - всё, завершаем рабочий день и отдых, уходим заниматься домашними делами
				setStepDone(taskId, step)
				setTaskDone(taskId)
				setPersonState(personId, "chores")
				log.Printf("%d | Персонаж %d закончил работать, сонливость - %.2f\n", nowTime, personId, somnolency)
			}
		} else {
			// А вот если спать ещё не хочется, смотрим на то, чем занимались
			// Если работали и при этом
			if taskType == "work" {
				// устали - отдыхаем
				if fatifue > conf.MAX_FATIGUE_FOR_STOP_WORK {
					setStepDone(taskId, step)
					setPersonState(personId, "rest")
					createNewStep("rest", step+1, taskId, finish_time)
					log.Printf("%d | Персонаж %d устал и решил отдохнуть, усталость - %.2f\n", nowTime, personId, fatifue)
				} else {
					// не устали - продолжаем работать
					setStepDone(taskId, step)
					createNewStep("work", step+1, taskId, finish_time)
				}
			} else if taskType == "rest" {
				// А вот если отдыхали и при этом
				// нормально отдохнули - продолжаем работать
				if fatifue < conf.MIN_FATIGUE_FOR_START_WORK {
					setStepDone(taskId, step)
					setPersonState(personId, "work")
					createNewStep("work", step+1, taskId, finish_time)
					log.Printf("%d | Персонаж %d отдохнул и продолжил работу, усталость - %.2f\n", nowTime, personId, fatifue)
				} else {
					// не успели отдохнуть - продолжаем отдых
					setStepDone(taskId, step)
					createNewStep("rest", step+1, taskId, finish_time)
				}
			}
		}
	}
}

func createNewStep(stepType string, step, taskId, start_time int) {
	//Получаем время завершения следующего шага
	stepFinishTime := getNextStepData(taskId, step)
	nowTime := lib.GetNowWorldTime()
	//Создаём новую задачу на работу
	// TODO В будущем нужно будет смотреть на тип задачи и для рабочих добавить результаты
	//fmt.Println("createNewStep")
	_, err := queries["createNewStep"].query.Exec(taskId, step, start_time, stepFinishTime, stepType, nowTime)
	//fmt.Println("!")
	if err != nil {
		log.Fatal(queries["createNewStep"].queryErrorText, step, taskId, err)
	}
}

func setStepDone(taskId, step int) {
	//fmt.Println("setStepDone")
	_, err := queries["setStepDone"].query.Exec(taskId, step)
	//fmt.Println("!")
	if err != nil {
		log.Fatal(queries["setStepDone"].queryErrorText, taskId, step, err)
	}
}

func setTaskDone(taskId int) {
	//fmt.Println("setTaskDone")
	_, err := queries["setTaskDone"].query.Exec(taskId)
	//fmt.Println("!")
	if err != nil {
		log.Fatal(queries["setTaskDone"].queryErrorText, taskId, err)
	}
}

func setPersonState(personId int, state string) {
	//fmt.Println(queries["setPersonState"].query, state, personId)
	_, err := queries["setPersonState"].query.Exec(state, personId)
	//fmt.Println("!")
	if err != nil {
		log.Fatalf(queries["setPersonState"].queryErrorText, state, personId, err)
	}
}

func getNextStepData(taskId, step int) int64 {
	var stepFinishTime sql.NullInt64
	nowTime := lib.GetNowWorldTime()
	for {
		//fmt.Println("getStepData")
		err := queries["getStepData"].query.QueryRow(taskId, step).Scan(&stepFinishTime)
		//fmt.Println("!")
		if err != nil {
			log.Fatalf(queries["getStepData"].queryErrorText, step, taskId, err)
		}
		// Нет нет запланированных шагов - уходим заниматься домашними делами
		if !stepFinishTime.Valid {
			// TODO Нужен вызов функции для перехода на занятие домашними делами
			return 0
		}
		// Если следующий шаг в будущем - принимем, иначе берём следующий шаг
		if stepFinishTime.Int64 > nowTime {
			return stepFinishTime.Int64
		}
	}
}
