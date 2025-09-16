package http

import (
	"encoding/json"
	"net/http"

	//"ecom-api/pkg/middleware"
	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/ports"

	"github.com/go-chi/chi/v5"
)

// Package http Payment Microservice API
//
// @title           Payment Microservice API
// @version         1.0
// @description     Handles payments for orders in the e-commerce system.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.email  support@example.com
//
// @host      localhost:8085
// @BasePath  /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type PaymentHandler struct {
	service ports.PaymentService
	orderClient ports.OrderClient
}

func NewPaymentHandler(s ports.PaymentService, oc ports.OrderClient) *PaymentHandler {
	return &PaymentHandler{service: s, orderClient: oc}
}




// @Summary      Create Payment
// @Description  Create a payment for a given order (requires Bearer token).
// @Tags         Payments
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "The ID of the order to pay for"
// @Success      200  {object}  domain.Payment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payments/{order_id} [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	
	userID, _ := middleware.FromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error": "unauthorized: missing or invalid token"}`, http.StatusUnauthorized)
		return
	}

	orderID := chi.URLParam(r, "order_id")
	if orderID == "" {
		http.Error(w, `{"error": "order ID is required"}`, http.StatusBadRequest)
		return
	}

	// 1. Fetch order from Order-MS
	order, err := h.orderClient.GetOrder(r.Context(), orderID)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch order: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// 2. Build payment object from order data
	payment := &domain.Payment{
		OrderID: order.Id,
		UserID:  order.UserId,
		Amount:  order.Total,
		Status:  "PENDING",
	}

	// 3. Persist via PaymentService
	created, err := h.service.ProcessPayment(r.Context(), payment)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
		
}
