package domain

type Product struct {
	ID   string  `json:"id" bson:"_id,omitempty"`
	Name  string  `json:"name"`
	Description   string  `json:"description"`
	Price   float64   `json:"price"`
	Stock   int     `json:"stock"`
}