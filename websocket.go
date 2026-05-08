package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
	"sync"
)

const wsChannelBuf = 64

type conn struct {
	stringc chan string
	bytec   chan []byte
}

type WSServer struct {
	mu        sync.Mutex
	chanList  *list.List
	device    *deviceMessage
	onControl func(controlMessage)
}

func NewWSServer(device *deviceMessage, onControl func(controlMessage)) *WSServer {
	return &WSServer{
		chanList:  list.New(),
		device:    device,
		onControl: onControl,
	}
}

func (s *WSServer) BroadcastMessage(data string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for e := s.chanList.Front(); e != nil; e = e.Next() {
		var c = e.Value.(conn)
		select {
		case c.stringc <- data:
		default:
		}
	}
}

func (s *WSServer) BroadcastBMessage(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for e := s.chanList.Front(); e != nil; e = e.Next() {
		var c = e.Value.(conn)
		select {
		case c.bytec <- data:
		default:
		}
	}
}

func (s *WSServer) closeN(c *list.Element) {
	s.mu.Lock()
	s.chanList.Remove(c)
	s.mu.Unlock()
}

func (s *WSServer) HandleMessages(c *websocket.Conn) {
	var cChannel = make(chan string, wsChannelBuf)
	var bChannel = make(chan []byte, wsChannelBuf)
	s.mu.Lock()
	var li = s.chanList.PushBack(conn{
		stringc: cChannel,
		bytec:   bChannel,
	})
	s.mu.Unlock()
	defer s.closeN(li)

	log.Println("New connection from", c.RemoteAddr())
	m, err := json.Marshal(s.device)
	if err != nil {
		log.Println("Error serializing JSON: ", err)
	}

	err = c.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
		return
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				return
			}
			var ctrl controlMessage
			if err := json.Unmarshal(msg, &ctrl); err != nil {
				continue
			}
			if ctrl.MessageType == "control" && s.onControl != nil {
				s.onControl(ctrl)
			}
		}
	}()

	running := true
	for running {
		select {
		case msg := <-cChannel:
			err = c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
				running = false
			}
		case msg := <-bChannel:
			err = c.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
				running = false
			}
		case <-done:
			running = false
		}
		runtime.Gosched()
	}
}

func (s *WSServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	s.HandleMessages(c)
}
