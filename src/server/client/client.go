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

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type User struct {
	hub *Hub

	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var (
	db      = conf.Db_client
	hub     = newHub()
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (u *User) readPump() {
	defer func() {
		u.hub.unregister <- u
		u.ws.Close()
	}()
	u.ws.SetReadLimit(maxMessageSize)
	u.ws.SetReadDeadline(time.Now().Add(pongWait))
	u.ws.SetPongHandler(func(string) error {
		u.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := u.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Fatal("error: %v", err)
			}
			break
		}
		fmt.Print(message)
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//u.hub.broadcast <- message
	}
}

func (u *User) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.ws.Close()
	}()
	for {
		select {
		case message, ok := <-u.send:
			u.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				u.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := u.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			u.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
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
	user := &User{hub: hub, ws: ws, send: make(chan []byte)}
	user.hub.register <- user
	go user.writePump()
	getInit(user.send)
	getPerson(user.send)
	user.readPump()
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/static/favicon.ico")
}

func ClientStart() {
	go hub.run()
	go getTime(hub.broadcast)

	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/client/static/", http.StripPrefix("/client/static/", http.FileServer(http.Dir("./client/static/"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(conf.ADDR, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
