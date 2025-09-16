package domain

//import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID         string   `json:"id" bson:"_id,omitempty"`
	UserID     string      `json:"user_id" bson:"user_id"`
	Items      []OrderItem `json:"items" bson:"items"`
	Total float64     `json:"total" bson:"total"`
	Status     string      `json:"status" bson:"status"` // PENDING, CONFIRMED, CANCELLED
}

type OrderItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	Price     float64 `json:"price" bson:"price"` // price per item
}
