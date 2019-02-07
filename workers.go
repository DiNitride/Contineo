//package main
//
//import (
//	"io/ioutil"
//	"fmt"
//	"encoding/json"
//	"github.com/gorilla/websocket"
//)
//
///*
//A worker is represented by the worker struct which contains all the information related to that
//specific worker node.
//*/
//
//func RegisterWorkers() {
//
//	workerFiles, err := ioutil.ReadDir("./workers/")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	for _, f := range workerFiles {
//		fmt.Printf("Loading worker file %v\n", f.Name())
//		data, err := ioutil.ReadFile("workers/" + f.Name())
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		var w Worker
//		var i id
//		err = json.Unmarshal(data, &w)
//		err = json.Unmarshal(data, &i)
//		if err != nil {
//			fmt.Println("Error decoding worker configuration")
//		} else {
//			Workers[i.Id] = &w
//		}
//	}
//}
//
//func PrintWorkers() {
//	for id, worker := range Workers {
//		fmt.Printf("ID: %v Game: %v Gamemode: %v\n", id, worker.Game, worker.Gamemode)
//	}
//}
//
//func ConnectNode(conn *websocket.Conn, data map[string]interface{}) {
//	/*
//	This method registers a connection to a loaded worker node.
//	 */
//	for nodeID, worker := range Workers {
//		if worker.Key == data["key"] {
//			// Valid origin, accept this as the worker node
//			if worker.Live {
//				// Worker already registered
//				fmt.Println("Worker already registered")
//				// TODO: Return error message to worker node
//				return
//			} else {
//				Workers[nodeID].Conn = conn
//				Workers[nodeID].Data = make(map[string]interface{})
//
//				resp := map[string]interface{}{
//					"origin": "controller",
//					"type": "ack",
//				}
//
//				jsonData, err := json.Marshal(resp);
//				if err != nil {
//					fmt.Printf("Error registering worker %v\n", conn.RemoteAddr())
//				}
//				conn.WriteMessage(1, []byte(jsonData))
//
//				fmt.Println("Node connected")
//			}
//		}
//	}
//}
//
//func GetWorkerIDFromConn(conn *websocket.Conn) string {
//	for nodeID, worker := range Workers {
//		if worker.Conn == conn {
//			return nodeID
//		}
//	}
//	return ""
//}
//
//func ProcessUpdate(conn *websocket.Conn, update map[string]interface{}) {
//	fmt.Println("Processing update data...")
//	NodeID := GetWorkerIDFromConn(conn)
//	for k, v := range update {
//		Workers[NodeID].Data[k] = v
//	}
//	fmt.Printf("Updated data for worker %v\n", conn.RemoteAddr())
//	Updates <- update
//}


package main

import "github.com/gorilla/websocket"

type Worker struct {
	connection *websocket.Conn
	accessToken string
	workerID string
	data map[string]interface{}
	errorCount int
}