package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Init orders var as a slice Order struct
var orders []Order
var dbName = "packerform-db"

// Route Handler for returning Orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// log.Println(r)
	json.NewEncoder(w).Encode(orders)
}

func main() {
	log.Println("Starting Golang Local Server.")
	// Init Router
	r := mux.NewRouter()

	// Prepare Data
	orders = prepareOrdersData()
	log.Println("Orders Data Ready.")

	// End points
	r.HandleFunc("/orders", getOrders).Methods(http.MethodGet, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))
	log.Println("Endpoints Ready.")

	// Start server and log if it fails
	// log.Fatal(http.ListenAndServeTLS(":8000", "./localhost.crt", "./localhost.key", nil))
	log.Println("Listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func prepareOrdersData() []Order {
	var returnOrders []Order
	mongoClient := getMongoClient()

	ordersCollection := getMongoCollection("orders", dbName, mongoClient)
	customersCollection := getMongoCollection("customers", dbName, mongoClient)
	companiesCollection := getMongoCollection("company", dbName, mongoClient)

	// set the orders
	orderCursor := getCollectionCursor(ordersCollection, bson.M{})
	for orderCursor.Next(context.TODO()) {
		var order Order
		err := orderCursor.Decode(&order)
		if err != nil {
			log.Fatal("orderCursor.Decode ERROR:", err)
		}
		returnOrders = append(returnOrders, order)
	}

	// Loop through orders
	for idx := range returnOrders {

		// Set the customer
		var customer Customer
		custID := returnOrders[idx].CustomerID
		customerfilter := bson.M{"user_id": bson.M{"$eq": custID}}
		setDecodedValue(customerfilter, customersCollection, &customer)

		// Set the company in customer
		companyFilter := bson.M{"company_id": bson.M{"$eq": customer.CompanyID}}
		setDecodedValue(companyFilter, companiesCollection, &customer.Company)

		// Set the order items
		var orderItems []OrderItem
		orderID := returnOrders[idx].ID
		orderItemsFilter := bson.M{"order_id": bson.M{"$eq": orderID}}
		setOrderItemsValue(orderItemsFilter, mongoClient, &orderItems)

		// Assign the customer to the order
		returnOrders[idx].OrderItems = orderItems
		returnOrders[idx].Customer = customer
	}

	return returnOrders
}

func setOrderItemsValue(filter bson.M, client *mongo.Client, orderItems *[]OrderItem) {
	var orderItem OrderItem

	orderItemsCollection := getMongoCollection("order-items", dbName, client)
	cursor := getCollectionCursor(orderItemsCollection, filter)

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&orderItem)
		if err != nil {
			log.Fatal("collection.Find ERROR:", err)
		}
		itemQuantity, _ := strconv.ParseFloat(orderItem.Quantity, 64)
		priceUnit, _ := strconv.ParseFloat(orderItem.PricePerUnit, 64)
		orderItem.OrderItemAmount = itemQuantity * priceUnit

		// For each order item get the delivery details
		deliveryFilter := bson.M{"order_item_id": bson.M{"$eq": orderItem.ID}}
		setDeliveriesValues(deliveryFilter, client, &orderItem, priceUnit)

		*orderItems = append(*orderItems, orderItem)
	}
}

func setDeliveriesValues(filter bson.M, client *mongo.Client, orderItem *OrderItem, priceUnit float64) {
	var deliveries []Delivery
	var delivery Delivery

	deliveryCollection := getMongoCollection("deliveries", dbName, client)
	cursor := getCollectionCursor(deliveryCollection, filter)

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&delivery)
		if err != nil {
			log.Fatal("collection.Find ERROR:", err)
		}
		deliveredQuanity, _ := strconv.ParseFloat(delivery.DeliveredQuantity, 64)
		delivery.DeliveredAmount = deliveredQuanity * priceUnit
		deliveries = append(deliveries, delivery)

	}

	orderItem.Deliveries = deliveries
}

func setDecodedValue(filter bson.M, collection *mongo.Collection, structObj interface{}) interface{} {
	cursor := getCollectionCursor(collection, filter)
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(structObj)
		if err != nil {
			log.Fatal("collection.Find ERROR:", err)
		}
	}

	return structObj
}

func getCollectionCursor(collection *mongo.Collection, filter bson.M) *mongo.Cursor {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return cursor
}
