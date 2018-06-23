package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var Workers = map[string]*Worker{}

func RegisterWorkers() {
	workerFiles, err := ioutil.ReadDir("./workers/")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range workerFiles {
		fmt.Printf("Loading worker file %v\n", f.Name())
		data, err := ioutil.ReadFile("workers/" + f.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		var w Worker
		var i id
		err = json.Unmarshal(data, &w)
		err = json.Unmarshal(data, &i)
		if err != nil {
			fmt.Println("Error decoding worker configuration")
		} else {
			Workers[i.Id] = &w
		}
	}
}

func PrintWorkers() {
	for id, worker := range Workers {
		fmt.Printf("ID: %v Game: %v Gamemode: %v\n", id, worker.Game, worker.Gamemode)
	}
}

// Connects a connection to a registered worker node
func ConnectNode(conn *websocket.Conn, data map[string]interface{}) {
	for nodeID, worker := range Workers {
		if worker.Key == data["key"] {
			// Valid origin, accept this as the worker node
			if worker.Live {
				// Worker already registered
				fmt.Println("Worker already registered")
				// TODO: Return error message to worker node
				return
			} else {
				Workers[nodeID].Conn = conn

				resp := map[string]interface{}{
					"origin": "controller",
					"type": "ack",
				}
				jsonData, err := json.Marshal(resp);
				if err != nil {
					fmt.Printf("Error registering worker %v\n", conn.RemoteAddr())
				}
				conn.WriteMessage(1, []byte(jsonData))

				fmt.Println("Node connected")
			}


		}
	}
}

func GetWorkerIDFromConn(conn *websocket.Conn) string {
	for nodeID, worker := range Workers {
		if worker.Conn == conn {
			return nodeID
		}
	}
	return ""
}

func ProcessUpdate(conn *websocket.Conn, update map[string]interface{}) {
	fmt.Println("Processing update data...")
	NodeID := GetWorkerIDFromConn(conn)
	for k, v := range update {
		Workers[NodeID].Data[k] = v
	}
	fmt.Printf("Updated data for worker %v\n", conn.RemoteAddr())
	Updates <- update
}

type Worker struct {
	Game string
	Gamemode string
	Description string
	Key string
	Data map[string]interface{}
	Live bool
	Conn *websocket.Conn
}

type id struct {
	Id string
}
