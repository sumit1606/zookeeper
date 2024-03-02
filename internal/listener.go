package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	conn, _, err := zk.Connect([]string{"localhost:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	nodePath := "/live_nodes"
	_, _, ch, err := conn.ChildrenW("/live_nodes")
	if err != nil {
		panic(err)
	}

	fmt.Println("Watching changes to node:", nodePath)

	for {
		select {
		case event := <-ch:
			if event.Type == zk.EventNodeChildrenChanged {
				data, _, err := conn.Get(nodePath)
				if err != nil {
					fmt.Println("Error getting node data:", err)
				} else {
					fmt.Printf("Node data changed: %s\n", string(data))
				}
			} else if event.Type == zk.EventNodeDeleted {
				fmt.Println("Node deleted")
				return
			}
		}
	}
}
