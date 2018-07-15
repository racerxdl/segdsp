package main

import (
	"github.com/gorilla/websocket"
	"container/list"
	"sync"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var wsMutex = sync.Mutex{}
var websocketConnections = list.New()
var currDevice = DeviceMessage{}