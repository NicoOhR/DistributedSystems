package main

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

	ms "github.com/jepsen-io/maelstrom/demo/go"
)

var COUNT int = 0
var mu = sync.Mutex{}

func unique_id(node *ms.Node) int {
	COUNT++
	// count = 0
	id_string := node.ID()
	node_number, _ := strconv.Atoi(id_string[1:])
	len := len(node.NodeIDs())
	return COUNT*len + node_number
}

func main() {
	id_node := ms.NewNode()
	id_node.Handle("generate", func(msg ms.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		mu.Lock()
		body["id"] = unique_id(id_node)
		mu.Unlock()
		id_node.Reply(msg, body)
		return nil
	})

	if err := id_node.Run(); err != nil {
		log.Fatal(err)
	}
}
