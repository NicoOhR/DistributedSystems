package main

import (
	"encoding/json"
	"log"
	"sync"

	ms "github.com/jepsen-io/maelstrom/demo/go"
)

type node struct {
	n        *ms.Node
	mu       *sync.Mutex
	messages []float64
	topology map[string][]string
}

type t_msg struct {
	topology map[string][]string `json:"topology"`
}

func (node node) broadcast() error {
	connected_nodes := node.topology[node.n.ID()]
	for _, adjacent_node := range connected_nodes {
		go func() {
			if error := node.n.Send(adjacent_node, node.messages); error != nil {
				panic(error)
			}
		}()
	}
	return nil
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
		node.mu.Lock()
		node.messages = append(node.messages, body["message"].(float64))

		if err := node.broadcast(); err != nil {
			return err
		}

		node.n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
		node.mu.Unlock()
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
