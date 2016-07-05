package main

import (
	"fmt"
	"log"
	"time"
)

/*
`
CREATE TABLE test_char
(
    person_id INTEGER DEFAULT nextval('test_char_person_id_seq'::regclass) NOT NULL,
    health NUMERIC DEFAULT 100.0 NOT NULL,
    fatigue NUMERIC DEFAULT 0 NOT NULL,
    hunger NUMERIC DEFAULT 0 NOT NULL,
    thirst NUMERIC DEFAULT 0 NOT NULL,
    somnolency NUMERIC DEFAULT 0 NOT NULL,
    state INTEGER DEFAULT 1 NOT NULL,
    last_htfs_update INTEGER DEFAULT 0 NOT NULL
);
`
*/

func test_insert() {
	const (
		INSERT_QUERY = "INSERT INTO test_char (last_htfs_update) VALUES ($1);"
		MAX_COUNT    = 100000
	)
	insert_query, err := db.Prepare(INSERT_QUERY)
	if err != nil {
		log.Fatal("Query preparation error -->%v\n", err)
	}
	per := 0
	t1 := time.Now()
	for i := 0; i < MAX_COUNT; i++ {
		_, err = insert_query.Exec(nowWorldTime)
		if err != nil {
			log.Fatal("Query execution error -->%v\n", err)
		}
		if i%int(MAX_COUNT/100) == 0 {
			per++
			fmt.Println(per, "% выполнено")
		}
	}
	t2 := time.Since(t1)
	fmt.Printf("%v queries are executed for %v seconds (%v per second)\n",
		MAX_COUNT, t2.Seconds(), MAX_COUNT/t2.Seconds())
}

func test_update() {
	const (
		UPDATE_QUERY = `UPDATE test_char SET
							hunger = hunger + ($1 * (($5 - last_htfs_update) / 3600.0)),
							thirst = thirst + ($2 * (($5 - last_htfs_update) / 3600.0)),
							fatigue = GREATEST(0.0, fatigue + ($3 * (($5 - last_htfs_update) / 3600.0))),
							somnolency = GREATEST(0.0, somnolency + ($4 * (($5 - last_htfs_update) / 3600.0))),
							last_htfs_update = last_htfs_update
						WHERE state = 1;`
		MAX_COUNT = 100000
	)
	update_query, err := db.Prepare(UPDATE_QUERY)
	if err != nil {
		log.Fatal("Query preparation error -->%v\n", err)
	}
	fmt.Println("Поехали!!")
	t1 := time.Now()
	/*
		for i := 0; i < MAX_COUNT; i++ {
			_, err = insert_query.Exec(nowWorldTime)
			if err != nil {
				log.Fatal("Query execution error -->%v\n", err)
			}
			if i%int(MAX_COUNT/100) == 0 {
				per++
				fmt.Println(per, "% выполнено")
			}
		}*/
	_, err = update_query.Exec(2.0, 4.0, 12.0, -1.0, 600)
	if err != nil {
		log.Fatal("Query execution error -->%v\n", err)
	}
	t2 := time.Since(t1)
	fmt.Printf("Query is executed for %v seconds (%v per second)\n",
		t2.Seconds(), MAX_COUNT/t2.Seconds())
}
