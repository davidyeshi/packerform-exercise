package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Open a mongo connection and return mongo client instance
func getMongoClient() *mongo.Client {
	// Connect to Mongo change port number if different
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Mongo Client.")

	return client
}

// gets a mongo collection given the collection name and mongo client
func getMongoCollection(collectionName string, dbName string, client *mongo.Client) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}
