package main

import (
	"encoding/json"
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

func getWorldCalendarTime(worldTime int64) map[string]string {
	cTime := make(map[string]string)
	for key, _ := range calendarConfig {
		t := int(worldTime/int64(calendarConfig[key].TimeInSeconds)) + calendarConfig[key].MinValue
		if (calendarConfig[key].PeriodLabel == "minute" || calendarConfig[key].PeriodLabel == "hour") && t < 10 {
			cTime[calendarConfig[key].PeriodLabel] = "0" + strconv.Itoa(t)
		} else {
			cTime[calendarConfig[key].PeriodLabel] = strconv.Itoa(t)
		}
		worldTime -= int64((t - calendarConfig[key].MinValue) * calendarConfig[key].TimeInSeconds)
	}
	return cTime
}

func getWCTString(cTime map[string]string) string {
	return cTime["year"] + " год, " + cTime["month"] + " месяц, " +
		cTime["ten_day"] + " декада, " + cTime["day"] + " день, " +
		cTime["hour"] + ":" + cTime["minute"]
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
