package db

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var Client *mongo.Client = DBinstance()

func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")
		if err != nil {
    		log.Fatalf("Error loading .env file: %v", err)
		}	

	MongoDb := os.Getenv("MONGO_URI")
	fmt.Print(MongoDb)

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to MongoDB")
	log.Println("Connected to MongoDB")
	return client
}




func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("recipedb").Collection(collectionName)

	return collection
}