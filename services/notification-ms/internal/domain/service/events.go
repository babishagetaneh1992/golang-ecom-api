package service



type OrderCreatedEvent struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Amount  float64 `json:"amount"`
	Status  string `json:"status"`
}


type PaymentCreatedEvent struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Amount    float64 `json:"amount"`
	Status    string `json:"status"`
}




