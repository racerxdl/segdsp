package main

import (
	"github.com/gorilla/websocket"
	"github.com/racerxdl/radioserver/protocol"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var currDevice = deviceMessage{
	Gain: protocol.Invalid,
}
