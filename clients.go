package main

import "github.com/gorilla/websocket"

type Client struct {
	token string
}

var Clients = map[*websocket.Conn]Client{}
