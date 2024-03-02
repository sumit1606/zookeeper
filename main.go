package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

var path = "/live_nodes"

func watchChanges() {
	conn, _, err := zk.Connect([]string{"localhost:2181"}, time.Second) //* change the address to your zookeeper server address
	if err != nil {
		panic(err)
	}
	_, _, ch, err := conn.ExistsW(path)
	if err != nil {
		panic(err)
	}
	for {
		event := <-ch
		log.Printf("Got event: %v", event)
		_, _, ch, err = conn.ExistsW(path) // Reset the watch
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	conn, _, err := zk.Connect([]string{"localhost:2181"}, time.Second) //* change the address to your zookeeper server address
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	createLiveNodes(conn)
	watchChanges()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8090", nil)
}

func createLiveNodes(conn *zk.Conn) {

	// Define the path and data
	data := []byte("all the live servers") // the data you want to store

	// Check if the node exists
	exists, _, err := conn.Exists(path)
	if err != nil {
		panic(err)
	}

	// If the node exists, delete it
	if exists {
		err = conn.Delete(path, -1)
		if err != nil {
			panic(err)
		}
	}

	// Create a node
	_, err = conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
}
