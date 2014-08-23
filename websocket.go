package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/kdar/factorlog"
	"net/http"
	"time"
)

const (
	writeTimeout   = 10 * time.Second
	pongTimeout    = 60 * time.Second
	pingPeriod     = (pongTimeout * 9) / 10
	maxMessageSize = 125 //size of control frames. Anything else is not allowed
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hub struct {
	connections map[*connection]bool
	broadcast   chan *wsMsg
	register    chan *connection
	unregister  chan *connection
}

var h = hub{
	connections: make(map[*connection]bool),
	broadcast:   make(chan *wsMsg),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			jsonOut, _ := json.Marshal(m)
			for c := range h.connections {
				select {
				case c.send <- jsonOut:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *connection) readHandler() {
	// Once we receive an error from the host, we want to deregister and remove
	// them from the pool
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongTimeout))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongTimeout)); return nil })
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeTimeout))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writeHandler() {
	pinger := time.NewTicker(pingPeriod)
	defer func() {
		pinger.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-pinger.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Nope", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorln("unable to upgrade to websocket: ", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writeHandler()
	c.readHandler()
}

func websocketInit() {
	go h.run()
	http.HandleFunc("/ws", serveWs)
}
