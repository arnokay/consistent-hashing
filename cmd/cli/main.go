package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/arnokay/consistent-hashing/pkg/notabisect"
	"github.com/arnokay/consistent-hashing/pkg/tools"
)

var servers []Server = []Server{
	{
		Addr: "localhost:3000",
		Name: "A",
	},
	{
		Addr: "localhost:3001",
		Name: "B",
	},
	{
		Addr: "localhost:3002",
		Name: "C",
	},
	{
		Addr: "localhost:3003",
		Name: "D",
	},
	{
		Addr: "localhost:3004",
		Name: "E",
	},
	{
		Addr: "localhost:3005",
		Name: "F",
	},
	{
		Addr: "localhost:3006",
		Name: "J",
	},
	{
		Addr: "localhost:3007",
		Name: "L",
	},
}

func main() {
	ch := NewConsistentHashing(80, 10)

	for _, server := range servers {
    err := ch.AddServer(server)
    if err != nil {
      log.Fatal(err)
    }
	}

	for i := 0; i < 10; i++ {
		ch.UploadFile(fmt.Sprintf("file_%v.png", i))
	}
}

func hexToBigInt(hexBytes []byte) (*big.Int, error) {
	bigInt := new(big.Int)
	bigInt = bigInt.SetBytes(hexBytes)
	return bigInt, nil
}

func InputToKey(input string) (*big.Int, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		return new(big.Int), err
	}
	b := hasher.Sum(nil)
	hexRepresentation, err := hexToBigInt(b)
	if err != nil {
		return new(big.Int), err
	}
	// key := new(big.Int).Mod(hexRepresentation, new(big.Int).SetInt64(int64(amout)))
	return hexRepresentation, nil
}

type Server struct {
	Addr string
	Name string
}

func (s Server) uploadFile(filename string) {
	fmt.Printf("uploading file %v to server %v\n", filename, s.Name)
	time.Sleep(1 * time.Second)
}

type ConsistentHashing struct {
	keys          []*big.Int
	servers       []Server
	replicaAmount int
	serversLimit  int
}

func NewConsistentHashing(serversLimit, replicaAmount int) *ConsistentHashing {
	if replicaAmount == 0 {
		replicaAmount = 1
	}
	return &ConsistentHashing{
		replicaAmount: replicaAmount,
		serversLimit:  serversLimit,
	}
}

func (ch *ConsistentHashing) AddServer(server Server) error {
	if len(ch.servers) >= ch.serversLimit {
		return errors.New("exided the limit on a servers")
	}
	for i := 0; i < ch.replicaAmount; i++ {
		input := fmt.Sprintf("%v_%v", server.Addr, i)
		key, err := InputToKey(input)
		if err != nil {
			message := fmt.Sprintf("cannot transform input (%v) to key", input)
			return errors.New(message)
		}
		index := notabisect.BisectBigInt(ch.keys, key)
		ch.servers = tools.Insert(ch.servers, server, index)
		ch.keys = tools.Insert(ch.keys, key, index)
	}
	return nil
}

func (ch *ConsistentHashing) UploadFile(filename string) {
	key, err := InputToKey(filename)
	if err != nil {
		log.Fatal(err)
	}
	index := notabisect.BisectRightBigInt(ch.keys, key)
	if index == len(ch.servers) {
		index = 0
	}
	server := ch.servers[index]
	server.uploadFile(filename)
}
