package main

import (
    "code.google.com/p/go.net/websocket"
    "net/http"
    "encoding/json"
)

type subscriber struct {
    conn *websocket.Conn
    ch chan bool
    s bool
}

type wsMsg struct {
    Type string
    Content interface{}
}

type headlineMsg struct {
    Icon,
    Source,
    Link,
    Text string
}

type alertMsg struct {
    Text string
}
    

var (
    subscriptions  = make(chan subscriber)
    messages       = make(chan *wsMsg, 5)
)

//func sendHeadline

func socketHandler() {
    conns := make(map[*websocket.Conn]chan bool)
    for {
        select {
        case sub := <-subscriptions:
            if sub.s {
                conns[sub.conn] = sub.ch
            } else {
                delete(conns, sub.conn)
            }
        case message := <-messages:
            jsonOut, _ := json.Marshal(message)
            for conn, ch := range conns {
                if _, err := conn.Write(jsonOut); err != nil {
                    conn.Close()
                    ch <- false
                    close(ch)
                }
            }
        }
    }
}

func wsHandler(ws *websocket.Conn) {
    ch := make(chan bool)
    subscriptions <- subscriber{ws, ch, true}
    subscriptions <- subscriber{ws, ch, <-ch}
}

func init() {
    go socketHandler()
    http.Handle("/ws", websocket.Handler(wsHandler))
}
