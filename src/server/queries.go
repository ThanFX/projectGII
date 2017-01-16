package main

import (
	"database/sql"
	"log"
)

type Query struct {
	query          *sql.Stmt
	text           string
	queryErrorText string
}

var (
	queries map[string]Query
)

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
	queries["updateAllPersonsHTFS"] = Query{
		emptyStmp,
		`UPDATE persons SET
			hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
			thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
			fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
			somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
			last_htfs_update = $5
		WHERE state = $6;`,
		`Ошибка глобального HTFS-обновления персонажей в состоянии %s: %s`}
	// $1 - skillId
	queries["getSkillInfo"] = Query{
		emptyStmp,
		`SELECT tools, results FROM skills WHERE id = $1;`,
		`Ошибка получения навыка %d: %s`}

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
