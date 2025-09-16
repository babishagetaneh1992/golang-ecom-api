package domain



type Payment struct {
	ID       string  `json:"id" bson:"_id,omitempty"`
	OrderID  string  `json:"order_id" bson:"order_id"`
	UserID   string  `json:"user_id" bson:"user_id"`
	Amount   float64 `json:"amount" bson:"amount"`
	Status   string  `json:"status" bson:"status"` // PENDING, COMPLETED, FAILED
	
	
	
}
