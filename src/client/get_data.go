package client

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ThanFX/projectGII/src/conf"
	"github.com/ThanFX/projectGII/src/lib"
)

type Chunk struct {
	id         int
	x          int
	y          int
	isExplored bool
	terrains   string
}

const (
	startMapX = 3
	startMapY = -5
	mapWidth  = 5
	mapHeight = 5
)

func getTime(sendTime chan []byte) {
	var world_time_str string
	var world_time int64
	for {
		err := db.QueryRow("SELECT world_time FROM time WHERE id = 1;").Scan(&world_time_str)
		if err != nil {
			log.Fatal("Ошибка запроса таймеров в БД: ", err)
		}
		world_time, err = strconv.ParseInt(world_time_str, 10, 64)
		//fmt.Println("Cчитанное из БД время: " + lib.GetWCTString(lib.GetWorldCalendarTime(world_time)))
		sendTime <- []byte("{\"key\":\"time\",\"value\":" + (lib.GetWCTJSON(lib.GetWorldCalendarTime(world_time))) + "}")
		//fmt.Println(lib.GetWCTJSON(lib.GetWorldCalendarTime(world_time)))
		time.Sleep(time.Second)
	}
}

func getInit(send chan []byte) {
	getConfig(send)
	getMap(startMapX, startMapY, send)
}

func getMap(startX, startY int, send chan []byte) {
	var chunk Chunk
	var out string
	worldMap, err := db.Query("SELECT * FROM world_map WHERE x >= $1 AND x < $2 AND world_map.y <= $3 AND world_map.y > $4;",
		startX, startX+mapWidth, startY, startY-mapHeight)
	if err != nil {
		log.Fatal("Ошибка запроса карты в БД: ", err)
	}
	defer worldMap.Close()
	out = "["
	for worldMap.Next() {
		err = worldMap.Scan(&chunk.id, &chunk.x, &chunk.y, &chunk.isExplored, &chunk.terrains)
		if err != nil {
			log.Fatal("Ошибка парсинга карты: ", err)
		}
		terrains := chunk.terrains[1 : len(chunk.terrains)-1]
		out += ("{\"x\":" + strconv.Itoa(chunk.x) + ",\"y\":" + strconv.Itoa(chunk.y) + ",\"isExplored\":" + strconv.FormatBool(chunk.isExplored) + "," + terrains + "},")
	}
	out = strings.TrimRight(out, ",")
	out += "]}"
	//fmt.Println(out)
	send <- []byte("{\"key\":\"worldMap\",\"value\":" + out)
}

func getPerson(send chan []byte) {
	var person conf.Person
	var out string
	for {
		persons, err := db.Query("SELECT p.id, p.name, p.chunk->'x', p.chunk->'y', p.state, p.health, p.hunger, p.thirst, p.fatigue, p.somnolency FROM persons p;")
		if err != nil {
			log.Fatal("Ошибка запроса пользователей в БД: ", err)
		}
		defer persons.Close()
		out = "["
		for persons.Next() {
			err = persons.Scan(&person.PersonId, &person.Name, &person.Chunk.X, &person.Chunk.Y, &person.State,
				&person.Health, &person.Hunger, &person.Thirst, &person.Fatigue, &person.Somnolency)
			if err != nil {
				log.Fatal("Ошибка парсинга списка персонажей: ", err)
			}
			personString, err := json.Marshal(person)
			if err != nil {
				log.Fatal("Ошибка преобразование персонажа в JSON: ", err)
			}
			out += (string(personString) + ",")
		}
		out = strings.TrimRight(out, ",")
		out += "]}"
		send <- []byte("{\"key\":\"persons\",\"value\":" + out)
		//fmt.Println(out)
		time.Sleep(time.Second * 5)
	}
}

func getConfig(send chan []byte) {
	mapInfo := "{\"startMapX\":" + strconv.Itoa(startMapX) + ",\"startMapY\":" + strconv.Itoa(startMapY) +
		",\"mapWidth\":" + strconv.Itoa(mapWidth) + ",\"mapHeight\":" + strconv.Itoa(mapHeight) + "}"
	send <- []byte("{\"key\":\"mapInfo\",\"value\":" + mapInfo + "}")
}
