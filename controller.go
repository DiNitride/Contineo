package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// A Channel that contains changes to servers waiting to be pushed to clients
var WorkerUpdates chan map[string]interface{}

// Slices to hold the current connections
var Workers []*Worker
var Clients []*Client

const ErrorLimit = 50

/*
Declaring errors
*/
type ContineoError int

const (
	jsonDecodeError ContineoError = iota
	unauthorised
)

func main() {
	WorkerUpdates = make(chan map[string]interface{})
	go ReadWorkerUpdate()
	http.HandleFunc("/controller", handleConnection)

	fmt.Println("Starting server. . .")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting http server")
	}

}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	/*
		Method that handles all new connections.
	*/

	fmt.Println("Handling new connection, authorising...")

	var accessToken string

	// First, handle authorisation
	connectionType := r.Header.Get("conn_type") // Client or worker
	if connectionType == "client" {
		token := r.Header.Get("token")
		/*
			Clients must be authorised against the external auth api, and the api must return an access token
			and a user id: this id must match the ID used within the workers to specify permissions.
			Workers are authorised against the loaded worker data
		*/
		resp, err := http.PostForm("http://127.0.0.1:5000/authenticate", url.Values{"token": {fmt.Sprintf("%s", token)}})
		if err != nil {
			fmt.Println("Error making post request to authentication server!")
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading authentication response")
		}
		resp.Body.Close()

		var authResponse map[string]interface{}
		err = json.Unmarshal(body, &authResponse)
		if err != nil {
			fmt.Println("Error json decoding authentication server response")
		}

		fmt.Println("Connection authorised.")
		accessToken = authResponse["access_token"].(string)
		if accessToken == "unauthorised" {
			fmt.Println("Invalid athentication token!")
			return
		}
		fmt.Printf("Access token is: %s\n", accessToken)
	} else if connectionType == "worker" {
		// Yeet
		fmt.Println("Worker!")
		accessToken = "worker_test"
	} else {
		fmt.Println("Invalid connection type specified in HTTP headers.")
		return
	}

	/*
		Upgrade to web socket connection.
		At this point, we want to store the connection somewhere we can easily access it later.
	*/
	fmt.Println("Upgrading connection to websocket")
	conn, err := upgrader.Upgrade(w, r, http.Header{"access_token": {accessToken}})
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
	err = json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println("Error decoding connections initialisation packet, aborting")
		return
	}
	// We now have the connections initialisation data

	fmt.Print("Initialisation data: ")
	fmt.Println(data)

	if accessToken != data["init_token"] {
		fmt.Println("Token mismatch between auth server and handshake, disconnecting!")
		conn.Close()
		return
	}

	if connectionType == data["origin"] {
		switch data["origin"] {
		case "worker":
			fmt.Println("Registering worker")
			RegisterWorker(conn, &data)
			break
		case "client":
			fmt.Println("Registering client")
			RegisterClient(conn, &data)
		default:
			// Error and Disconnect
		}
	} else {
		fmt.Println("Connection type mismatch between handshake packet and original HTTP headers")
	}
}

func RegisterWorker(conn *websocket.Conn, data *map[string]interface{}) {
	/*
		Method to register a new worker node with the controller
	*/
	if (*data)["type"] == "init" {
		worker := Worker{connection: conn, data: *data}
		Workers = append(Workers, &worker)
		sendRecvAcknowledgment(worker.connection)
		go RecvWorkerData(&worker)
		fmt.Println("New worker node connected")
	} else {
		fmt.Println("Error registering new worker, did not send init packet.")
	}
}

func RegisterClient(conn *websocket.Conn, data *map[string]interface{}) {
	/*
		Method to register a new client node with the controller
	*/
	if (*data)["type"] == "init" {
		client := Client{connection: conn}
		Clients = append(Clients, &client)
		sendRecvAcknowledgment(client.connection)
		fmt.Println("New client node connected")
	} else {
		fmt.Println("Error registering new worker, did not send init packet.")
	}
}

func sendRecvAcknowledgment(conn *websocket.Conn) {
	/*
		Method to send an acknowledgement packet to a connections
	*/
	resp := map[string]interface{}{
		"origin": "controller",
		"type":   "ack",
	}

	jsonData, err := json.Marshal(resp)
	err = (*conn).WriteMessage(1, []byte(jsonData))
	if err != nil {
		fmt.Printf("Error sending ack packet to %v\n", (*conn).RemoteAddr())
	}
}

func sendError(conn *websocket.Conn, code ContineoError) {
	/*
		Method to send an error packet to a connection.
	*/
	resp := map[string]interface{}{
		"origin": "controller",
		"type":   "error",
		"code":   code,
	}

	jsonData, err := json.Marshal(resp)
	err = (*conn).WriteMessage(1, []byte(jsonData))
	if err != nil {
		fmt.Printf("Error sending ack packet to %v\n", (*conn).RemoteAddr())
	}
}

func RecvWorkerData(worker *Worker) {
	/*
		Method to loop and pull data from the worker connections. A new goroutine is spawned for each connection.
	*/

	// Interface object to load json data into
	var data map[string]interface{}

	for {

		// Check the error count

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

		json.Unmarshal(msg, &data)
		/*
			"Unmarshal" means convert from json string to object. It takes the string and a pointer to the above string
			interface object thing.

			We now have JSON Data and can start dealing with it :D
		*/

		// Send an ack to the client
		sendRecvAcknowledgment((*worker).connection)

		WorkerUpdates <- data
		fmt.Printf("Recieved data from %v\n", (*worker).connection.RemoteAddr())
	}
}

func ReadWorkerUpdate() {
	/*
		This method continuously reads data from the Updates channel and processes it.
	*/
	for {
		update := <-WorkerUpdates
		fmt.Println(update)
		Broadcast(update)
	}
}

func Broadcast(update map[string]interface{}) {

	jsonData, err := json.Marshal(update)
	if err != nil {
		fmt.Println("Error Marshalling update")
	}
	for _, c := range Clients {
		c.connection.WriteMessage(1, []byte(jsonData))
	}
}

var upgrader = websocket.Upgrader{
	/*
		Upgrades a HTTP request to /controller to a web socket connection.
	*/
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
