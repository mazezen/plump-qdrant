package main

import (
	"os"

	"github.com/mazezen/plum-qdrant/go-qdrant/collections"
	"github.com/mazezen/plum-qdrant/go-qdrant/retrieve"
	"github.com/qdrant/go-client/qdrant"
)

// go run main.go [argument]
func main() {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "",
		Port: 6334,
		APIKey: "",
		UseTLS: true,
	})
	if err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) < 2 {
		panic("缺少参数")
	}
	p := args[1]
	switch p {
	case "create": 
		collections.CreateCollection(client)
	case "get": 
		collections.GetCollectionDetails(client)
	case "check": 
		collections.CheckCollectionExists(client)
	case "update": 
		collections.UpdateCollectionParameters(client)
	case "del": 
		collections.DeleteCollectonByGiven(client)
	case "list":
		collections.ListAllCollections(client)
	case "upsert":
		retrieve.NewRetrieve(client).UpsertPoints()
	}

	
	// collections.CreateCollection(client)
	// collections.GetCollectionDetails(client)
	// collections.UpdateCollectionParameters(client)
	// collections.CheckCollectionExists(client)
	// collections.DeleteCollectonByGiven(client)
}