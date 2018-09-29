package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
)

var chanList = list.New()

type conn struct {
	stringc chan string
	bytec   chan []byte
}

func closeN(c *list.Element) {
	wsMutex.Lock()
	chanList.Remove(c)
	wsMutex.Unlock()
}

func broadcastMessage(data string) {
	wsMutex.Lock()
	for e := chanList.Front(); e != nil; {
		var c = e.Value.(conn)
		go func() {
			c.stringc <- data
		}()
		var next = e.Next()
		e = next
	}
	wsMutex.Unlock()
}
func broadcastBMessage(data []byte) {
	wsMutex.Lock()
	for e := chanList.Front(); e != nil; {
		var c = e.Value.(conn)
		go func() {
			c.bytec <- data
		}()
		var next = e.Next()
		e = next
	}
	wsMutex.Unlock()
}

func handleMessages(c *websocket.Conn) {

	var cChannel = make(chan string)
	var bChannel = make(chan []byte)
	wsMutex.Lock()
	var li = chanList.PushBack(conn{
		stringc: cChannel,
		bytec:   bChannel,
	})
	wsMutex.Unlock()
	defer closeN(li)

	// region Send DeviceInfo
	log.Println("New connection from", c.RemoteAddr())
	m, err := json.Marshal(currDevice)
	if err != nil {
		log.Println("Error serializing JSON: ", err)
	}

	err = c.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
		return
	}
	// endregion
	// region Client Loop
	running := true
	for running {
		//_, _, err := c.ReadMessage()
		//if err != nil {
		//	break
		//}
		select {
		case msg := <-cChannel:
			err = c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
				running = false
				break
			}
		case msg := <-bChannel:
			err = c.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
				running = false
				break
			}
		}

		runtime.Gosched()
	}
	// endregion
}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	handleMessages(c)
}
