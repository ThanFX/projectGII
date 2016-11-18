package client

import (
	"fmt"
	"log"
	"net/http"
	"server/conf"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	err := homeTemplate.Execute(w, r.Host)
	if err != nil {
		log.Fatal(w, "Template error: ", err)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	fmt.Println("Соединение установлено!")
}

func ClientStart() {
	go getTime()

	r := mux.NewRouter()
	http.HandleFunc("/", homeHandler)
	r.HandleFunc("/ws", wsHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(conf.ADDR, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
