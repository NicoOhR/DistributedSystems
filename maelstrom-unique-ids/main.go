package main

import (
	"encoding/json"
	"log"
	"sync"

	ms "github.com/jepsen-io/maelstrom/demo/go"
)

var COUNT int = -1
var mu = sync.Mutex{}

func unique_id() int {
	COUNT++
	return COUNT
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
		body["id"] = unique_id() * len(id_node.NodeIDs())
		mu.Unlock()
		id_node.Reply(msg, body)
		return nil
	})

	if err := id_node.Run(); err != nil {
		log.Fatal(err)
	}
}
