package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

type Node struct {
	node *snowflake.Node
	mu   sync.Mutex
}

func NewNode(nodeID int64) *Node {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		panic(err)
	}
	return &Node{
		node: node,
	}
}

func (n *Node) GenerateID() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.node.Generate().String()
}

func (n *Node) GenerateIDInt64() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.node.Generate().Int64()
}
