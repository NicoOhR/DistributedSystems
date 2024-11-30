package main

import (
	"encoding/json"
	"log"

	ms "github.com/jepsen-io/maelstrom/demo/go"
)

type node struct {
	n        *ms.Node
	messages []float64
	topology map[string][]string
}

type t_msg struct {
	topology map[string][]string `json:"topology"`
}

func main() {

	node := node{
		n:        ms.NewNode(),
		messages: []float64{},
	}

	node.n.Handle("broadcast", func(msg ms.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		node.messages = append(node.messages, body["message"].(float64))
		node.n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
		return nil
	})

	node.n.Handle("read", func(msg ms.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "read_ok"
		body["messages"] = node.messages
		node.n.Reply(msg, body)
		return nil
	})

	node.n.Handle("topology", func(msg ms.Message) error {
		var message t_msg
		body := make(map[string]any)
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			return err
		}

		node.topology = message.topology

		body["type"] = "topology_ok"
		node.n.Reply(msg, body)
		return nil
	})

	if err := node.n.Run(); err != nil {
		log.Fatal(err)
	}
}
