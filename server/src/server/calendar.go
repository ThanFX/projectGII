package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/patrickmn/sortutil"
)

type CalendarPeriod struct {
	MaxValue      int    `json: "maxValue"`
	MinValue      int    `json: "minValue"`
	PeriodName    string `json: "periodName"`
	PeriodLabel   string `json: "periodLabel"`
	TimeInSeconds int    `json: "timeInSeconds"`
}

type CalendarTime struct {
	Period string
	Time   string
}

var (
	calendarConfig []CalendarPeriod
)

func getWorldCalendarTime(worldTime int64) []CalendarTime {
	cTime1 := make(map[string]string)
	cTime := make([]CalendarTime, len(calendarConfig))
	for key, _ := range calendarConfig {
		t := int(worldTime/int64(calendarConfig[key].TimeInSeconds)) + calendarConfig[key].MinValue
		cTime[key].Period = calendarConfig[key].PeriodLabel
		cTime[key].Time = strconv.Itoa(t)
		cTime1[calendarConfig[key].PeriodLabel] = strconv.Itoa(t)
		worldTime -= int64((t - calendarConfig[key].MinValue) * calendarConfig[key].TimeInSeconds)
	}
	fmt.Println(cTime1)
	return cTime
}

func getWCTString(cTime []CalendarTime) string {
	return cTime[0].Time + " год, " +
		cTime[1].Time + " месяц, " +
		cTime[2].Time + " декада, " +
		cTime[3].Time + " день, " +
		cTime[4].Time + ":" +
		cTime[5].Time
}

func getCalendar() {
	var res string
	err := db.QueryRow("SELECT value->'periods' FROM config WHERE id = 'calendar';").Scan(&res)
	if err != nil {
		log.Fatal("Ошибка запроса календаря в БД", err)
	}
	bytes := []byte(res)
	err = json.Unmarshal(bytes, &calendarConfig)
	if err != nil {
		log.Fatal("Ошибка парсинга структуры календаря", err)
	}
	sortutil.DescByField(calendarConfig, "TimeInSeconds")
}
