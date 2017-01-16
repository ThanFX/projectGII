package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"server/conf"
	"server/lib"
	"strconv"
	"strings"
	"time"
)

// Создаём задачи (и массив шагов выполнения для задач) на старт работ для персонажей
func create_task() {
	var personId, taskId int
	var skill conf.Skill
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
		skill = getSkillInfo(preferSkillId)
		//fmt.Println(skill.Results.ItemTemplateId)
		steps := getTaskSteps(nowTime, skill.Results.ItemTemplateId)
		//fmt.Println(steps)
		//fmt.Println("createTask")
		err = queries["createTask"].query.QueryRow(personId, preferSkillId, steps, "{}", nowTime).Scan(&taskId)
		//fmt.Println("!")
		if err != nil {
			log.Fatalf(queries["createTask"].queryErrorText, personId, startDayTime, err)
		}
		//createNewStep("wait", 0, taskId, int(nowTime))
		fmt.Println("Создали работу для ", personId)
	}
}

func getSkillInfo(preferSkillId int) conf.Skill {
	var skill conf.Skill
	//fmt.Println("getSkillInfo")
	err := queries["getSkillInfo"].query.QueryRow(preferSkillId).
		Scan(&skill.Tools, &skill.Results.ItemTemplateId, &skill.Results.Type)
	if err != nil {
		log.Fatalf(queries["getSkillInfo"].queryErrorText, preferSkillId, err)
	}
	return skill
}

func getTaskSteps(nowTime int64, itemParentId int) string {
	// Первый шаг через некоторый промежуток после пробуждения
	stepTime := rand.Int63n(conf.MAX_MORNING_INTERVAL_BEFORE_WORK) + conf.MIN_STEP_DURATING + nowTime
	steps := "{"
	step := 0
	// Получаем массив айдишников результатов
	items := getFoods(itemParentId)
	// Создаём массив шагов, до конца текущих суток, нулевой шаг - ошидание работы с момента пробуждения,
	// после него начинается реальная работа
	for stepTime < (nowTime + 24*3600) {
		steps += "\"" + strconv.Itoa(step) + "\":{\"finish_time\":" + strconv.FormatInt(stepTime, 10)
		// Запрашиваем возможный результат работы для данного шага
		res, itemNum, quantity := getStepFoodResult(len(items))
		steps += ",\"is_res\":" + strconv.FormatBool(res) + ",\"item_id\":" + strconv.Itoa(items[itemNum]) +
			",\"quantity\":" + strconv.Itoa(quantity)
		steps += "},"
		step++
		stepTime += rand.Int63n(conf.MAX_STEP_DURATING-conf.MIN_STEP_DURATING) + conf.MIN_STEP_DURATING
	}
	steps = strings.TrimRight(steps, ",")
	steps += "}"
	return steps
}

func getStepFoodResult(maxItems int) (bool, int, int) {
	var itemNum, quantity int
	res := false
	resultChance := rand.Intn(100)
	//fmt.Print("Шанс результата - ", resultChance, " | ")
	if resultChance >= 20 {
		res = true
		itemNum = rand.Intn(maxItems)
		quantity = rand.Intn(4) + 1
	}
	return res, itemNum, quantity
}

func getFoods(foodParentId int) []int {
	var food int
	var foods []int

	results, err := queries["getFoods"].query.Query(foodParentId)
	if err != nil {
		log.Fatalf(queries["getFoods"].queryErrorText, foodParentId, err)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&food)
		if err != nil {
			log.Fatalf("Ошибка парсинга шаблона продукта: %s", err)
		}
		foods = append(foods, food)
	}
	return foods
}

func getPreferPersonSkill(personId int) int {
	var skillId int
	//fmt.Println("getPreferSkill)
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
