package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"encoding/json"
)

// A Channel that contains changes to servers waiting to be pushed to clients
var Updates chan map[string]interface{}


func main() {
	Updates = make(chan map[string]interface{})
	RegisterWorkers()
	PrintWorkers()
	go ReadUpdate()
	http.HandleFunc("/controller", handleConnection)

	fmt.Println("Starting server. . .")
	http.ListenAndServe(":3000",nil)

}

func handleConnection(w http.ResponseWriter, r *http.Request) {

	// First, handle authorisation
	token_type := r.Header.Get("token_type")
	token := r.Header.Get("token")




	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("New connection %v\n", conn.RemoteAddr())

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var data map[string]interface{}
		json.Unmarshal(msg, &data)
		fmt.Println(data)
		fmt.Printf("Recieved data from %v\n", conn.RemoteAddr())
		// We now have JSON Data and can start dealing with it

		if data["origin"] == "worker" {
			// Data from a worker
			// Most likely an update packet
			switch data["type"] {
			case "init":
				ConnectNode(conn, data)
			case "update":
				ProcessUpdate(conn, data)
			}
		}
		//else if data["origin"] == "client" {
		//	switch data["type"] {
		//	case "init":
		//		registerClient(conn)
		//	}
		//}

	}
}

func ReadUpdate() {
	for {
		update := <- Updates
		fmt.Println(update)
	}
}

//func registerClient(conn *websocket.Conn) {
//	clients = append(clients, conn)
//	resp := map[string]interface{}{
//		"origin": "controller",
//		"type": "ack",
//	}
//	jsonData, err := json.Marshal(resp);
//	if err != nil {
//		fmt.Printf("Error registering client %v\n", conn.RemoteAddr())
//	}
//	conn.WriteMessage(1, []byte(jsonData))
//	fmt.Printf("Registered client %v\n", conn.RemoteAddr())
//}
//
//func sendUpdates() {
//	for {
//		update := <- updates
//		jsonData, err := json.Marshal(update);
//		if err != nil {
//			fmt.Println(err)
//		}
//		for _, client := range clients {
//			for _, worker := range workers {
//				if worker.ready {
//
//					client.WriteMessage(1, []byte(jsonData))
//				}
//			}
//		}
//	}
//}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}