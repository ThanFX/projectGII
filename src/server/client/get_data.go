package client

import (
	"log"
	"server/lib"
	"strconv"
	"strings"
	"time"
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
	startMapY = -9
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
		time.Sleep(time.Second)
	}
}

func getInit(send chan []byte) {
	getMap(startMapX, startMapY, send)
}

func getMap(startX, startY int, send chan []byte) {
	var chunk Chunk
	var out string
	worldMap, err := db.Query("SELECT * FROM world_map WHERE x >= $1 AND x < $2 AND world_map.y >= $3 AND world_map.y < $4;",
		startX, startX+mapWidth, startY, startY+mapHeight)
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
	send <- []byte("{\"key\":\"worldMap\",\"value\":" + out)
}
