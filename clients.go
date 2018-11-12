package main

import "github.com/gorilla/websocket"

type Client struct {
	connection *websocket.Conn
	clientID string
	accessToken string
}
