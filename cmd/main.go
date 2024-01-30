package main

import (
	"fmt"
	"log"

	conhash "github.com/arnokay/consistent-hashing/internal/consistent-hashing"
)

func main() {
	chConfig := conhash.NewConfig().WithServersLimit(150000).WithReplicaAmount(10)
	ch := conhash.New(chConfig)

	nodesAmount := 15000
	for i := 0; i < nodesAmount; i++ {
		addr := fmt.Sprintf("localhost:300%v", i)
		err := ch.AddNode(conhash.NewNode(addr))
		if err != nil {
			log.Fatal(err)
		}
	}

	filesAmount := 10
	for i := 0; i < filesAmount; i++ {
		filename := fmt.Sprintf("file%v.txt", i)
		err := ch.UploadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}
