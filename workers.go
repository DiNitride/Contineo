package main

import "github.com/gorilla/websocket"

type Worker struct {
	connection *websocket.Conn
	workerID   string
	data       map[string]interface{}
}
