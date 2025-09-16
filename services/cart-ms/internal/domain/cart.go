package domain

import "time"

type CartItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Name      string  `json:"name" bson:"name"`
	Price     float64 `json:"price" bson:"price"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	AddedAt   time.Time `json:"added_at" bson:"added_at"`
}

type Cart struct {
	UserID    string     `json:"user_id" bson:"user_id"`
	Items     []CartItem `json:"items" bson:"items"`
	Total     float64    `json:"total" bson:"-"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
}
