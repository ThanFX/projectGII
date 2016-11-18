package client

import (
	"server/conf"
)

var (
	db = conf.Db_client
)

func ClientStart() {
	go getTime()
}
