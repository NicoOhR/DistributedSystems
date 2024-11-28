package main

import (
	"encoding/json"
	"log"

	ms "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var ids []int

	echo_node := ms.NewNode()
	echo_node.Handle("echo", func(msg ms.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "echo_ok"

		return echo_node.Reply(msg, body)
	})

	if err := echo_node.Run(); err != nil {
		log.Fatal(err)
	}

}
