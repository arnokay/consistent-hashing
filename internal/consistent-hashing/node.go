package consistenthashing

import (
	"fmt"
	"time"
)

type Node struct {
	key  string
	addr string
}

func NewNode(addr string) Node {
	return Node{
		key:  "",
		addr: addr,
	}
}

func (n Node) uploadFile(filename string) {
	fmt.Printf("uploading file %v to server %v\n", filename, n.addr)
	time.Sleep(1 * time.Second)
}
