package main

import (
	"log"
	"time"

	"github.com/amovah/eznode"
	"github.com/amovah/eznodefiber"
	"github.com/sirupsen/logrus"
)

func main() {
	node1 := eznode.NewChainNode(eznode.NewChainNodeConfig{
		Name: "node 1",
		Url:  "https://example.com",
		Limit: eznode.ChainNodeLimit{
			Count: 10,
			Per:   5 * time.Second,
		},
		RequestTimeout: 10 * time.Second,
		Priority:       1,
		Middleware:     nil, // optional
	})

	node2 := eznode.NewChainNode(eznode.NewChainNodeConfig{
		Name: "node 2",
		Url:  "https://example.com",
		Limit: eznode.ChainNodeLimit{
			Count: 10,
			Per:   5 * time.Second,
		},
		RequestTimeout: 10 * time.Second,
		Priority:       2,
		Middleware:     nil, // optional
	})

	chain := eznode.NewChain(eznode.NewChainConfig{
		Id: "Ethereum",
		Nodes: []*eznode.ChainNode{
			node1,
			node2,
		},
		CheckTickRate: eznode.CheckTick{
			TickRate:         100 * time.Millisecond,
			MaxCheckDuration: 5 * time.Second,
		},
	})

	createdEzNode := eznode.NewEzNode([]*eznode.Chain{chain})

	err := eznodefiber.StartFiber(8080, createdEzNode, logrus.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}
}
