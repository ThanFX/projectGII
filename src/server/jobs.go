package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	queries     map[string]*Query
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
	queries = make(map[string]*Query)
	// $1 - nowTime
	queries["getPersonWithoutWork"].text = `SELECT person_id FROM tasks WHERE daytime > $1;`
	queries["getPersonWithoutWork"].queryErrorText = `Ошибка выполнения запроса получения списка персонажей без задач на сегодня: %s`
	// $1 - person_id
	queries["getPreferSkill"].text = `SELECT skill_id FROM person_skills WHERE person_id = $1 ORDER BY worth DESC LIMIT 1;`
	queries["getPreferSkill"].queryErrorText = `Ошибка получения наилучшей работы для персонажа %d: %s`
	//
	queries["createTask"].text = `INSERT INTO tasks(person_id, skill_id, daytime, steps, result, create_time)
    	VALUES ($1, $2, $3, $4, $5, $6)`
	queries["createTask"].queryErrorText = `Ошибка создания задачи на работу для персонажа %d в %d: $s`
	// $1 - state, $2 - personId
	queries["setPersonState"].text = `UPDATE person_health_characteristic SET state = $1 WHERE person_id = $2;`
	queries["setPersonState"].queryErrorText = `Ошибка установки состояния \"%s\" для персонажа %d: %s`
	//$1 - taskId, $2 - step
	queries["setStepDone"].text = `UPDATE task_steps SET is_done = TRUE WHERE task_id = $1 AND step = $2;`
	queries["setStepDone"].queryErrorText = `Ошибка закрытия шага %d для задачи %d: %s`
	//$1 - taskId, $2 - step
	queries["getStepData"].text = `SELECT (steps->$2)->'finish_time' from tasks WHERE id = $1;`
	queries["getStepData"].queryErrorText = `Ошибка получения данных шага %d задачи %d: %s`
	//$1 - taskId
	queries["setTaskDone"].text = `UPDATE tasks SET is_done = TRUE WHERE task_id = $1;`
	queries["setTaskDone"].queryErrorText = `Ошибка закрытия задачи %d: %s`
	//
	queries["createNewStep"].text = `INSERT INTO task_steps(task_id, step, start_time, finish_time, type, create_time)
    	VALUES ($1, $2, $3, $4, $5, $6);`
	queries["createNewStep"].queryErrorText = `Ошибка создания шага %d задачи %d: %s`
	queries["getFinishStep"].text = `SELECT ts.task_id, t.person_id, t.skill_id, t.type, t.step, t.finish_time, chr.fatigue, chr.somnolency FROM tasks t
		JOIN person_health_characteristic chr ON chr.person_id = t.person_id
		WHERE t.type IN ('work', 'rest') AND t.is_done = FALSE AND t.finish_time < $1;`
}

func prepareQueries() {
	initQueries()
	for key, _ := range queries {
		prepareQuery, err := db.Prepare(queries[key].text)
		if err != nil {
			log.Fatal("Ошибка подготовки запроса ", key, ": ", err)
		}
		queries[key].query = prepareQuery
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
	state_period := strconv.Itoa(int((check.State * 60) / world_time_speed))
	state_period = "@every " + state_period + "s"
	hts_period := strconv.Itoa(int((check.HTS * 60) / world_time_speed))
	hts_period = "@every " + hts_period + "s"
	ed_period := strconv.Itoa(int((check.ED * 60) / world_time_speed))
	ed_period = "@every " + ed_period + "s"
	work_period := strconv.Itoa(int((check.Work * 60) / world_time_speed))
	work_period = "@every " + ed_period + "s"

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

	task_cron := cron.New()
	task_cron.AddFunc(work_period, task_job)
	go task_job()
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
		UPDATE person_health_characteristic SET state = $1 WHERE
			somnolency > 40.0;
		`, "sleep")
		if err != nil {
			log.Fatal("Ошибка засыпания персонажей: ", err)
		}
	}
	if cH >= 6 && cH < 20 {
		_, err := db.Exec(`
		UPDATE person_health_characteristic SET state = $1 WHERE
			somnolency <= 6.0;
		`, "chores")
		if err != nil {
			log.Fatal("Ошибка пробуждения персонажей: ", err)
		}
		// И сразу после пробуждения начали создавать задачи на сегодняшнюю работу
		go create_task()
	}

}

func ed_job() {
	if lib.GetNowWorldTime() == 0 {
		return
	}

	//cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	//fmt.Println("Едим и пьём в ", lib.GetWCTString(cTime))

	_, err := db.Exec(`
		UPDATE person_health_characteristic SET thirst = $1 WHERE
			thirst >= 6.0 AND state != $2;
		`, 0.0, "sleep")
	if err != nil {
		log.Fatal("Ошибка утоления жажды: ", err)
	}

	_, err = db.Exec(`
		UPDATE person_health_characteristic SET thirst = $1, hunger = $2 WHERE
			hunger >= 3.0 AND state != $3;
		`, 0.0, 0.0, "sleep")
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
		UPDATE person_health_characteristic SET
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
		UPDATE person_health_characteristic SET
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
		UPDATE person_health_characteristic SET
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
		UPDATE person_health_characteristic SET
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
	var personId int
	nowTime := lib.GetNowWorldTime()
	startDayTime := lib.GetStartDayTime(nowTime)
	// Получаем список персонажей, у которых нет работы
	persons, err := queries["getPersonWithoutWork"].query.Query(nowTime)
	if err != nil {
		log.Fatal(queries["getPersonWithoutWork"].queryErrorText, err)
	}
	defer persons.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	for persons.Next() {
		err = persons.Scan(&personId)
		if err != nil {
			log.Fatal("Ошибка парсинга списка персонажей: ", err)
		}
		preferSkillId := getPreferPersonSkill(personId)
		nowTime := lib.GetNowWorldTime()
		steps := getTaskSteps(nowTime)
		_, err = queries["createTask"].query.Exec(personId, preferSkillId, startDayTime, steps, "{}", nowTime)
		if err != nil {
			log.Fatalf(queries["createTask"].queryErrorText, personId, startDayTime, err)
		}
		fmt.Println("Создали работу для ", personId)
	}
}

func getTaskSteps(nowTime int64) string {
	// Первый шаг через некоторый промежуток после пробуждения
	stepTime := rand.Int63n(conf.MAX_MORNING_INTERVAL_BEFORE_WORK) + conf.MIN_STEP_DURATING + nowTime
	steps := "{"
	step := 1
	// Создаём массив шагов, до конца текущих суток
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
	err := queries["getPreferSkill"].query.QueryRow(personId).Scan(&skillId)
	if err != nil {
		log.Fatalf(queries["getPreferSkill"].queryErrorText, personId, err)
	}
	return skillId
}

func task_job() {
	go step_task_job()
}

/*
func create_task_job() {
	var taskId, personId, skillId, stepFinishTime int
	//Ищем открытые таски на создание работы, которые уже начались
	taskCreateQuery, _ := db.Prepare(`SELECT id, person_id, skill_id FROM tasks
		WHERE is_done = FALSE AND type = 'create' AND finish_time < $1;`)
	firstStepFinishQuery, _ := db.Prepare(`SELECT (steps->$1)->'finish_time' from task_steps
		WHERE task_id = $2;`)

	nowTime := lib.GetNowWorldTime()
	taskCreateList, err := taskCreateQuery.Query(nowTime)
	if err != nil {
		log.Fatal("Ошибка получения наступивших задач на старт работ: ", err)
	}
	defer taskCreateList.Close()
	// Для каждой такой таски запрашиваем первый шаг, создаём таску на работу по нему и переводим состояние персонажа в "work"
	for taskCreateList.Next() {
		err = taskCreateList.Scan(&taskId, &personId, &skillId)
		if err != nil {
			log.Fatal("Ошибка парсинга списка наступивших задач на старт работ: ", err)
		}
		//Получаем время завершения первого шага
		step := 1
		err = firstStepFinishQuery.QueryRow(step, taskId).Scan(&stepFinishTime)
		if err != nil {
			log.Fatal("Ошибка получения времени завершения первого шага работы ", taskId, ": ", err)
		}
		// Если шаг уже завершен - перебираем шаги, пока не наткнёмся на завершенный
		for int64(stepFinishTime) < nowTime {
			step++
			err = firstStepFinishQuery.QueryRow(step, taskId).Scan(&stepFinishTime)
			if err != nil {
				log.Fatal("Ошибка получения времени завершения ", step, " шага работы ",
					taskId, ": ", err)
			}
		}
		// Оборачиваем следующие изменения в транзакцию
		tx, err := db.Begin()
		if err != nil {
			log.Fatal("Ошибка открытия транзакции создания задачи на работу: ", err)
		}
		defer tx.Rollback()
		createJobQuery, _ := tx.Prepare(`INSERT INTO tasks(person_id, skill_id, start_time, finish_time, type,
			result, create_time, step) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)
		defer createJobQuery.Close()
		taskCreateDoneQuery, _ := tx.Prepare(`UPDATE tasks SET is_done = TRUE WHERE id = $1;`)
		defer taskCreateDoneQuery.Close()
		personWorkStateQuery, _ := tx.Prepare(`UPDATE person_health_characteristic SET state = 4
			WHERE person_id = $1;`)
		defer personWorkStateQuery.Close()

		//Закрываем задачу на создание работы
		_, err = taskCreateDoneQuery.Exec(taskId)
		if err != nil {
			log.Fatal("Ошибка закрытия задачи на старт работы ", taskId, ": ", err)
		}
		//Переводим персонажа в "рабочий режим"
		_, err = personWorkStateQuery.Exec(personId)
		if err != nil {
			log.Fatal("Ошибка перевода персонажа ", personId, " в состояние \"работа\": ", err)
		}
		// Создаём задачу на первый шаг работы
		_, err = createJobQuery.Exec(personId, skillId, nowTime, stepFinishTime, "work", "{}", nowTime, step)
		if err != nil {
			log.Fatal("Ошибка создания задачи на ", step, " шаг работы ", taskId, ": ", err)
		}
		err = tx.Commit()
		if err != nil {
			log.Fatal("Ошибка открытия транзакции создания задачи на работу: ", err)
		}
		log.Println("Успешно стартанули работу ", taskId, " для персонажа ", personId)
	}
}
*/

func step_task_job() {
	var taskId, personId, step, skillId, finish_time int
	var fatifue, somnolency float32
	var taskType string
	tasksQuery, _ := db.Prepare(`SELECT t.id, t.person_id, t.skill_id, t.type, t.step, t.finish_time, chr.fatigue, chr.somnolency FROM tasks t
		JOIN person_health_characteristic chr ON chr.person_id = t.person_id
		WHERE t.type IN ('work', 'rest') AND t.is_done = FALSE AND t.finish_time < $1;`)

	nowTime := lib.GetNowWorldTime()
	tasksList, err := tasksQuery.Query(nowTime)
	if err != nil {
		log.Fatal("Ошибка получения наступивших задач на шаги работы: ", err)
	}
	defer tasksList.Close()
	for tasksList.Next() {
		err = tasksList.Scan(&taskId, &personId, &skillId, &taskType, &step, &finish_time, &fatifue, &somnolency)
		if err != nil {
			log.Fatal("Ошибка парсинга списка наступивших задач на шаги работы: ", err)
		}
		fmt.Println("смотрим персонажа ", personId, " задача ", taskId)
		// Обрабатываем окончание работы - смотрим на повышенную сонливость
		if somnolency > conf.MAX_SOMNOLENCY_FOR_STOP_WORK {
			// Если не устали и при этом работали - продолжаем пахать
			if fatifue < conf.MAX_FATIGUE_FOR_STOP_WORK && taskType == "work" {
				setTaskDone(taskId)
				newStep("work", step+1, taskId, skillId, personId, finish_time)
			} else {
				// Если же устали или уже отдыхали на работе - всё, завершаем рабочий день и отдых, уходим заниматься домашними делами
				setTaskDone(taskId)
				setPersonState(personId, "chores")
			}
		} else {
			// А вот если спать ещё не хочется, смотрим на то, чем занимались
			// Если работали и при этом
			if taskType == "work" {
				// устали - отдыхаем
				if fatifue > conf.MAX_FATIGUE_FOR_STOP_WORK {
					setTaskDone(taskId)
					setPersonState(personId, "rest")
					newStep("rest", step+1, taskId, skillId, personId, finish_time)
				} else {
					// не устали - продолжаем работать
					setTaskDone(taskId)
					newStep("work", step+1, taskId, skillId, personId, finish_time)
				}
			} else if taskType == "rest" {
				// А вот если отдыхали и при этом
				// нормально отдохнули - продолжаем работать
				if fatifue < conf.MIN_FATIGUE_FOR_START_WORK {
					setTaskDone(taskId)
					setPersonState(personId, "work")
					newStep("work", step+1, taskId, skillId, personId, finish_time)
				} else {
					// не успели отдохнуть - продолжаем отдых
					setTaskDone(taskId)
					newStep("rest", step+1, taskId, skillId, personId, finish_time)
				}
			}
		}
	}
}

func createNewStep(stepType string, step, taskId, start_time int) {
	//Получаем время завершения следующего шага
	stepFinishTime := getNextStepData(taskId, step)
	nowTime := lib.GetNowWorldTime()
	//Создаём новую задачу на работу (на следующий шаг) и закрываем предыдущую
	// TODO В будущем нужно будет смотреть на тип задачи и для рабочих добавить результаты
	_, err := queries["createNewStep"].query.Exec(taskId, step, start_time, stepFinishTime, stepType, nowTime)
	if err != nil {
		log.Fatal(queries["createNewStep"].queryErrorText, step, taskId, err)
	}
}

func setStepDone(taskId, step int) {
	_, err := queries["setStepDone"].query.Exec(taskId, step)
	if err != nil {
		log.Fatal(queries["setStepDone"].queryErrorText, taskId, step, err)
	}
}

func setTaskDone(taskId int) {
	_, err := queries["setTaskDone"].query.Exec(taskId)
	if err != nil {
		log.Fatal(queries["setTaskDone"].queryErrorText, taskId, err)
	}
}

func setPersonState(personId int, state string) {
	_, err := queries["setPersonState"].query.Exec(state, personId)
	if err != nil {
		log.Fatalf(queries["setPersonState"].queryErrorText, state, personId, err)
	}
}

func getNextStepData(taskId, step int) int64 {
	var stepFinishTime sql.NullInt64
	nowTime := lib.GetNowWorldTime()
	for {
		err := queries["getStepData"].query.QueryRow(taskId, step).Scan(&stepFinishTime)
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
