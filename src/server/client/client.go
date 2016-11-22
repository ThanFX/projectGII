package client

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/conf"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	//hub *Hub

	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

var (
	db       = conf.Db_client
	sendTime chan string
)

func (p *Player) sendTime() {
	for {

		p.ws.WriteMessage(websocket.TextMessage, []byte(<-sendTime))
		fmt.Println("Test!")
		time.Sleep(time.Second)
	}
}

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

	homeTemplate, err := template.ParseFiles("client/static/templates/index.html")

	if err != nil {
		log.Fatal("Ошибка парсинга шаблона index.html: ", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	homeTemplate.ExecuteTemplate(w, "index", nil)
	/*
		if err != nil {
			log.Fatal(w, "Template error: ", err)
			return
		}
	*/
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	fmt.Println("Соединение установлено!")
	player := &Player{ws: ws, send: sendTime}

	go player.sendTime()
}

func ClientStart() {
	go getTime(sendTime)

	http.Handle("/client/static/", http.StripPrefix("/client/static/", http.FileServer(http.Dir("./client/static/"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(conf.ADDR, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
