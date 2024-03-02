package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func main() {
	register()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8082", nil)
}

func register() {
	// Connect to Zookeeper
	conn, _, err := zk.Connect([]string{"localhost:2181"}, time.Second) //* change the address to your zookeeper server address
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Define the path
	path := "/live_nodes/server2" // the path where you want to create your ephemeral node

	data := []byte("8081")
	// Create an ephemeral node
	_, err = conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
}
