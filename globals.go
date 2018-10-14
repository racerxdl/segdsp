package main

import (
	"github.com/gorilla/websocket"
	"github.com/racerxdl/spy2go/spyserver"
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
	Gain: spyserver.InvalidValue,
}
