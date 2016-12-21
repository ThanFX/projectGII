package main

import (
	"encoding/json"
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

	task_cron := cron.New()
	task_cron.AddFunc(state_period, task_job)
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
		create_works()
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

	//cTime := lib.GetWorldCalendarTime(lib.GetNowWorldTime())
	//fmt.Println("Проверка работы началась в ", lib.GetWCTString(cTime))

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

// Создаём задачи (и массив шагов выполнения для задач) на старт работ для персонажей
func create_works() {
	var personId, countTask, preferSkillId, createTaskId int
	// Получаем список персонажей, которые находятся в состоянии "домашние дела"
	persons, err := db.Query("SELECT p.id FROM persons p JOIN person_health_characteristic chr ON chr.person_id = p.id WHERE chr.state = 5;")
	if err != nil {
		log.Fatal("Ошибка запроса пользователей в БД: ", err)
	}
	defer persons.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	for persons.Next() {
		err = persons.Scan(&personId)
		if err != nil {
			log.Fatal("Ошибка парсинга списка персонажей: ", err)
		}
		// Получаем количество запланированных (время начала работы в будущем) работ по каждому персонажу
		err = db.QueryRow("SELECT count(id) FROM tasks WHERE person_id = $1 AND start_time > $2;", personId, lib.GetNowWorldTime()).Scan(&countTask)
		if err != nil {
			log.Fatal("Ошибка получения количества запланированных задач: ", err)
		}
		// Нет запланированных работ - назначаем
		if countTask == 0 {
			createTask, err := db.Prepare(`INSERT INTO tasks(person_id, skill_id, start_time, finish_time, type, result, create_time)
    			VALUES ($1, $2, $3, $3, $4, $5, $6) RETURNING id`)
			if err != nil {
				log.Fatal("Ошибка подготовки запроса на создание работы: ", err)
			}
			createTaskSteps, err := db.Prepare(`INSERT INTO task_steps(task_id, steps) VALUES ($1, $2)`)
			if err != nil {
				log.Fatal("Ошибка подготовки запроса на создание шагов работы: ", err)
			}
			err = db.QueryRow("SELECT skill_id FROM person_skills WHERE person_id = $1 ORDER BY worth DESC LIMIT 1;", personId).Scan(&preferSkillId)
			if err != nil {
				log.Fatal("Ошибка получения наилучшей работы: ", err)
			}
			randTime := rand.Int63n(6600) + 600
			nowTime := lib.GetNowWorldTime()

			err = createTask.QueryRow(personId, preferSkillId, nowTime+randTime, "create", "{}", nowTime).Scan(&createTaskId)
			if err != nil {
				log.Fatal("Ошибка при создании задачи на новую работу: ", err)
			}

			if err != nil {
				log.Fatal("Ошибка при получении id созданной задачи: ", err)
			}
			//fmt.Println("Создали работу ", createTaskId, " для ", personId)

			steps := "{"
			step := 1
			stepTime := nowTime + randTime
			for stepTime < (nowTime + 12*3600) {
				stepTime += rand.Int63n(3000) + 600
				steps += ("\"" + strconv.Itoa(step) + "\":{\"finish_time\":" + strconv.FormatInt(stepTime, 10))
				steps += "},"
				step++
			}
			steps = strings.TrimRight(steps, ",")
			steps += "}"
			//fmt.Println(steps)
			_, err = createTaskSteps.Exec(createTaskId, steps)
			if err != nil {
				log.Fatal("Ошибка при создании шагов новой работы: ", err)
			}
		}
	}
}

func task_job() {
	go create_task_job()
}

func create_task_job() {
	var taskId, personId, skillId, stepFinishTime int
	//Ищем открытые таски на создание работы, которые уже начались
	taskCreateQuery, _ := db.Prepare(`SELECT id, person_id, skill_id FROM tasks
		WHERE is_done = FALSE AND type = 'create' AND finish_time < $1;`)
	taskCreateDoneQuery, _ := db.Prepare(`UPDATE tasks SET is_done = TRUE WHERE id = $1;`)
	personWorkStateQuery, _ := db.Prepare(`UPDATE person_health_characteristic SET state = 4 WHERE person_id = $1;`)
	firstStepFinishQuery, _ := db.Prepare(`SELECT (steps->$1)->'finish_time' from task_steps WHERE task_id = $2;`)
	createJobQuery, _ := db.Prepare(`INSERT INTO tasks(person_id, skill_id, start_time, finish_time, type, result, create_time, step)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)

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
			log.Fatal("Ошибка получения времени заверешения первого шага работы ", taskId, ": ", err)
		}
		// Если шаг уже завершен - перебираем шаги, пока не наткнёмся на завершенный
		for int64(stepFinishTime) < nowTime {
			step++
			err = firstStepFinishQuery.QueryRow(step, taskId).Scan(&stepFinishTime)
			if err != nil {
				log.Fatal("Ошибка получения времени заверешения ", step, " шага работы ", taskId, ": ", err)
			}
		}
		// Создаём задачу на первый шаг работы
		_, err = createJobQuery.Exec(personId, skillId, nowTime, stepFinishTime, "work", "{}", nowTime, step)
		if err != nil {
			log.Fatal("Ошибка создания задачи на ", step, " шаг работы ", taskId, ": ", err)
		}
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
		log.Println("Успешно стартанули работу ", taskId, " для персонажа ", personId)
	}
}
