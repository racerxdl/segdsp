package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"container/list"
	"encoding/json"
	"runtime"
)

var chanList = list.New()

func closeN(c *list.Element) {
	wsMutex.Lock()
	chanList.Remove(c)
	wsMutex.Unlock()
}

func broadcastMessage(data string) {
	wsMutex.Lock()
	for e := chanList.Front(); e != nil; {
		var c = e.Value.(chan string)
		go func() {
			c <- data
		}()
		var next = e.Next()
		e = next
	}
	wsMutex.Unlock()
}

func handleMessages(c *websocket.Conn) {

	var cChannel = make(chan string)
	wsMutex.Lock()
	var li = chanList.PushBack(cChannel)
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
	for {
		//_, _, err := c.ReadMessage()
		//if err != nil {
		//	break
		//}
		msg := <- cChannel
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("Error sending message:", err, "dropping connection from", c.RemoteAddr())
			break
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