package http

import (
	"encoding/json"
	"net/http"

	//"order-microservice/internals/ports"

	//"ecom-api/pkg/middleware"

	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/ports"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service ports.OrderService
}

func NewOrderHandler(s ports.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}



// @Summary      Create Order
// @Description  Place a new order for the authenticated user from their cart
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [post]
// CreateOrder handles placing a new order from the user's cart
func (s *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.FromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error": "unauthorized: missing userID"}`, http.StatusUnauthorized)
		return
	}

	// Call service method to create order from cart
	createdOrder, err := s.service.CreateOrderFromCart(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order placed successfully",
		"order":   createdOrder,
	})
}

// GetOrder fetches a single order by ID
func (s *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, err := s.service.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// ListOrders returns all orders
func (s *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.service.ListOrders(r.Context())
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus updates the status of an order and returns the updated order
func (s *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, `{"error": "status field is required"}`, http.StatusBadRequest)
		return
	}

	updatedOrder, err := s.service.UpdateOrderStatus(r.Context(), id, req.Status)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "order status updated successfully",
		"order":   updatedOrder,
	})
}

// DeleteOrder deletes an order by ID
func (s *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := s.service.DeleteOrder(r.Context(), id); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Order deleted successfully"}`))
}
