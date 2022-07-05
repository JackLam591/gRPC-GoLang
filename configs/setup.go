package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	LoadEnv()
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGOURI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // setup 10 seconds timeout for connecting DB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

//Client instance
var DB *mongo.Client = ConnectDB()

//getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("GoLangOnboard").Collection(collectionName)
	return collection
}
