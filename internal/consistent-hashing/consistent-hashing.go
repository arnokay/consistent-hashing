package consistenthashing

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/arnokay/consistent-hashing/pkg/notabisect"
	"github.com/arnokay/consistent-hashing/pkg/tools"
)

var (
	defaultReplicaAmount = 1
	defaultServersLimit  = 0
)

type config struct {
	replicaAmount int
	serversLimit  int
}

func NewConfig() config {
	return config{
		replicaAmount: defaultReplicaAmount,
		serversLimit:  defaultServersLimit,
	}
}

func (c config) WithServersLimit(limit int) config {
	c.serversLimit = limit
	return c
}

func (c config) WithReplicaAmount(amount int) config {
	c.replicaAmount = amount
	return c
}

type ConsistentHashing struct {
	keys   []string
	nodes  []Node
	config config
}

func New(c config) *ConsistentHashing {
	return &ConsistentHashing{
		config: c,
	}
}

func (ch *ConsistentHashing) AddNode(n Node) error {
	if len(ch.nodes) > ch.config.serversLimit {
    message := fmt.Sprintf("exided the limit (%v) on a servers", ch.config.serversLimit)
    fmt.Printf("%v", ch.nodes)
		return errors.New(message)
	}
	for i := 0; i < ch.config.replicaAmount; i++ {
		input := fmt.Sprintf("%v_%v", n.addr, i)
		key, err := ch.hash(input)
		if err != nil {
			return err
		}
		index := notabisect.Bisect(ch.keys, key)
		ch.nodes = tools.Insert(ch.nodes, n, index)
		ch.keys = tools.Insert(ch.keys, key, index)
	}
	return nil
}

func (ch *ConsistentHashing) UploadFile(filename string) error {
	key, err := ch.hash(filename)
	if err != nil {
		return err
	}
	index := notabisect.BisectRight(ch.keys, key)
	if index == len(ch.keys) {
		index = 0
	}
	node := ch.nodes[index]
	node.uploadFile(filename)
	return nil
}

func (ch *ConsistentHashing) hash(input string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		message := fmt.Sprintf("cannot transform input (%v) to key", input)
		return "", errors.New(message)
	}
	b := hasher.Sum(nil)

	return string(b), nil
}
