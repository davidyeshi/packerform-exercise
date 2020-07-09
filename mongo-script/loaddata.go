package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// instantiating mongo client
	mongoClient := getMongoClient()
	fileAndCollectionNameMap := getFileAndCollection()
	dbName := "packerform-db"

	for collectionName, filePath := range fileAndCollectionNameMap {
		loadData(mongoClient, filePath, collectionName, dbName)
	}

	// Disconnecting connection
	mongoClient.Disconnect(context.TODO())
}

// Gets a map of collection and the data file path
func getFileAndCollection() map[string]string {
	return map[string]string{
		"company":     "./test_data/Test task - Mongo - customer_companies.csv",
		"customers":   "./test_data/Test task - Mongo - customers.csv",
		"orders":      "./test_data/Test task - Orders.csv",
		"order-items": "./test_data/Test task - Postgres - order_items.csv",
		"deliveries":  "./test_data/Test task - Postgres - deliveries.csv"}
}

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

	fmt.Println("Connected to Mongo Client.\n")

	return client
}

// Load the Companies into Mongo
func loadData(mongoClient *mongo.Client, filePath string, collectionName string, dbName string) {
	fmt.Printf("Loading " + collectionName + "...\n")

	csvLines := openCsvFile(filePath)

	// Loop through the lines and get the headers
	headers := csvLines[0]
	collection := getMongoCollection(collectionName, dbName, mongoClient)

	for _, line := range csvLines[1:] {
		toInsert := bson.M{}

		for lineIndex, header := range headers {
			toInsert[header] = line[lineIndex]
		}

		addDocToMongo(collection, toInsert)
	}
	fmt.Printf("Loaded " + strconv.Itoa(len(csvLines[1:])) + " doc(s) into " + collectionName + ".\n\n")
}

// Open a csv file and return its lines
func openCsvFile(filePath string) [][]string {
	csvFile, err := os.Open(filePath)
	defer csvFile.Close() // closes once function returns

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully opened ", filePath)
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	return csvLines
}

// gets a mongo collection given the collection name and mongo client
func getMongoCollection(collectionName string, dbName string, client *mongo.Client) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}

// Add a row/doc into a collection
func addDocToMongo(collection *mongo.Collection, toInsert bson.M) {
	_, err := collection.InsertOne(context.TODO(), toInsert)
	if err != nil {
		log.Fatal(err)
	}
}
