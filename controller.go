package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// A Channel that contains changes to servers waiting to be pushed to clients
var WorkerUpdates chan map[string]interface{}

// Slices to hold the current connections
var Workers []*Worker


func main() {
	WorkerUpdates = make(chan map[string]interface{})
	go ReadWorkerUpdate()
	http.HandleFunc("/controller", handleConnection)

	fmt.Println("Starting server. . .")
	http.ListenAndServe(":3000",nil)

}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	/*
	Method that handles all new connections.
	 */

	fmt.Println("Handling new connection, authorising...")
	// First, handle authorisation
	//token_type := r.Header.Get("token_type")
	//token := r.Header.Get("token")
	/*
	Clients must be authorised against the external auth api, and the api must return an access token
	and a user id: this id must match the ID used within the workers to specify permissions.
	Workers are authorised against the loaded worker data
	 */
	fmt.Println("Connection authorised.")
	accessToken := "token"

	// IF IT GET'S PAST HERE IT'S AUTHORISED YEET

	/*
	Upgrade to web socket connection.
	At this point, we want to store the connection somewhere we can easily access it later.
	 */
	fmt.Println("Upgrading connection to websocket")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("New connection %v\n", conn.RemoteAddr())

	// Here we must wait and pass the connection to a method to get it's initialisation data.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]interface{}
	json.Unmarshal(msg, &data)
	// We now have the connections initialisation data

	fmt.Print("Initialisation data: ")
	fmt.Println(data)

	switch data["origin"] {
	case "worker":
		RegisterWorker(conn, &data, accessToken)
	case "client":
		// RegisterClient
	default:
		// Error and Disconnect
	}


}

func RegisterWorker (conn *websocket.Conn, data *map[string]interface{}, token string, ) {
	/*
	Method to register a new worker node with the controller
	*/
	if (*data)["type"] == "init" {
		worker := Worker{connection:conn,data:*data,accessToken:token}
		Workers = append(Workers, &worker)
		sendAck(&worker)
		go RecvWorkerData(&worker)
		fmt.Println("New worker node connected")
	} else {
		fmt.Println("Error registering new worker, did not send init packet.")
	}
}

func sendAck(worker *Worker) {
	/*
	Method to send an acknowledgement packet to a connections
	*/
	resp := map[string]interface{}{
		"origin": "controller",
		"type": "ack",
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error sending ack packet to %v\n", (*worker). connection.RemoteAddr())
	}
	(*worker).connection.WriteMessage(1, []byte(jsonData))


}

func RecvWorkerData (worker *Worker) {
	/*
	Method to loop and pull data from the worker connections. A new goroutine is spawned for each connection.
	*/

	for {

		fmt.Printf("Waiting for data from... %v\n", (*worker).connection.RemoteAddr())

		_, msg, err := (*worker).connection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		/*
		Dealing with JSON is a bit of a pain. First you have to a map of string to interfaces. To be perfectly honest
		with you, not entirely sure what this means, but I know that it works and that's how I can access JSON data
		without having to specifically declare what values it contains.
		 */
		var data map[string]interface{}
		json.Unmarshal(msg, &data)
		/*
		"Unmarshal" means convert from json string to object. It takes the string and a pointer to the above string
		interface object thing.
		 */

		 // Send an ack to the client
		 sendAck(worker)

		// TODO: Ensure acccess token in packet matches stored token
		WorkerUpdates <- data
		fmt.Printf("Recieved data from %v\n", (*worker).connection.RemoteAddr())
		// We now have JSON Data and can start dealing with it :D

		// For now we are just accepting data from worker nodes
		if data["origin"] == "worker" {
			// Data from a worker
			// Most likely an update packet, so we'll check for that first
			switch data["type"] {
			case "update":
				// ProcessUpdate(conn, data)
			case "init":
				// ConnectNode(conn, data)
			}
		}
	}
}

func ReadWorkerUpdate() {
	/*
	This method continuously reads data from the Updates channel and processes it.
	*/
	for {
		update := <- WorkerUpdates
		fmt.Println(update)
	}
}

var upgrader = websocket.Upgrader{
	/*
	Upgrades a HTTP request to /controller to a web socket connection.
	 */
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}