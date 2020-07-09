package main

// Order Struct/Model
type Order struct {
	ID         string      `bson:"id" json:"order_id"`
	OrderDate  string      `bson:"created_at" json:"order_date"`
	OrderName  string      `bson:"order_name" json:"order_name"`
	CustomerID string      `bson:"customer_id" json:"-"`
	Customer   Customer    `json:"customer"`
	OrderItems []OrderItem `json:"order_items"`
}

// Customer Struct/Model
type Customer struct {
	CustID    string  `bson:"user_id" json:"customer_id"`
	CustName  string  `bson:"name" json:"customer_name"`
	CompanyID string  `bson:"company_id" json:"-"`
	Company   Company `json:"company"`
}

// Company Struct/Model
type Company struct {
	ID          string `bson:"company_id" json:"company_id"`
	CompanyName string `bson:"company_name" json:"company_name"`
}

// OrderItem Struct/Model
type OrderItem struct {
	ID              string     `bson:"id" json:"order_item_id"`
	OrderID         string     `bson:"order_id" json:"order_id"`
	PricePerUnit    string     `bson:"price_per_unit" json:"price_per_unit"`
	Quantity        string     `bson:"quantity" json:"quantity"`
	ProductName     string     `bson:"product" json:"product_name"`
	OrderItemAmount float64    `json:"order_item_amount"`
	Deliveries      []Delivery `json:"deliveries"`
}

// Delivery Struct/Model
type Delivery struct {
	ID                string  `bson:"id" json:"delivery_id"`
	DeliveredQuantity string  `bson:"delivered_quantity" json:"delivered_quantity"`
	DeliveredAmount   float64 `json:"delivery_amount"`
	OrderItemID       string  `bson:"order_item_id" json:"-"`
}
