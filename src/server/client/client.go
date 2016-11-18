package client

import (
	"log"
	"net/http"
	"server/conf"
	"text/template"
)

var (
	db           = conf.Db_client
	homeTemplate = template.Must(template.ParseFiles("client/public/index.html"))
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, r.Host)
}

func ClientStart() {
	go getTime()

	http.HandleFunc("/", homeHandler)
	err := http.ListenAndServe(conf.ADDR, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
