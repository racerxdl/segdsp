package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"container/list"
	"encoding/json"
)

func closeN(c *list.Element) {
	wsMutex.Lock()
	var v = c.Value.(*websocket.Conn)
	log.Println("Closing connection from", v.RemoteAddr())
	websocketConnections.Remove(c)
	c.Value.(*websocket.Conn).Close()
	wsMutex.Unlock()
}

func broadcastMessage(data string) {
	wsMutex.Lock()
	for e := websocketConnections.Front(); e != nil; {
		var ws = e.Value.(*websocket.Conn)
		err := ws.WriteMessage(websocket.TextMessage, []byte(data))
		var next = e.Next()
		if err != nil {
			log.Println("Error sending message:", err, "dropping connection from", ws.RemoteAddr())
			websocketConnections.Remove(e)
			ws.Close()
		}
		e = next
	}
	wsMutex.Unlock()
}

func handleMessages(c *websocket.Conn) {
	var el = websocketConnections.PushBack(c)
	defer closeN(el)

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
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
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