package main

import (
	"github.com/gorilla/websocket"
	"github.com/racerxdl/spy2go"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var wsMutex = sync.Mutex{}
var currDevice = deviceMessage{
	Gain: spy2go.InvalidValue,
}
