package main

import (
	"github.com/gorilla/websocket"
	"sync"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var wsMutex = sync.Mutex{}
var currDevice = DeviceMessage{}