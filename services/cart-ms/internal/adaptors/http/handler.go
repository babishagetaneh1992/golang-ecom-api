package http

import (
	//"ecom-api/pkg/middleware"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/ports"
	//"cart-microservice/internal/domain"
	//"cart-microservice/internal/ports"
)

// @title           Cart Microservice API
// @version         1.0
// @description     This is the Cart service for the e-commerce system.
// @host            localhost:8082
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
type CartHandler struct {
	service ports.CartService
}

type AddItemRequest struct {
    ProductID string `json:"product_id" validate:"required"`
    Quantity  int    `json:"quantity" validate:"required,min=1"`
}


func NewCartHandler(service ports.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// @Summary      Get Cart
// @Description  Get the authenticated user's cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /carts [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.FromContext(r.Context())
	if userID == "" {
		http.Error(w, "unauthorized: userID missing in context", http.StatusUnauthorized)
		return
	}

	cart, err := h.service.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type CartItemResponse struct {
		ProductID string  `json:"productId"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Quantity  int     `json:"quantity"`
		Subtotal  float64 `json:"subtotal"`
		AddedAt   string  `json:"addedAt"`
	}

	type CartResponse struct {
		UserID    string             `json:"userId"`
		Items     []CartItemResponse `json:"items"`
		Total     float64            `json:"total"`
		UpdatedAt string             `json:"updatedAt"`
	}

	var items []CartItemResponse
	var total float64
	for _, it := range cart.Items {
		subtotal := it.Price * float64(it.Quantity)
		items = append(items, CartItemResponse{
			ProductID: it.ProductID,
			Name:      it.Name,
			Price:     it.Price,
			Quantity:  it.Quantity,
			Subtotal:  subtotal,
			AddedAt:   it.AddedAt.Format("2006-01-02 15:04"),
		})
		total += subtotal
	}

	resp := CartResponse{
		UserID:    cart.UserID,
		Items:     items,
		Total:     total,
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02 15:04"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Add Item to Cart
// @Description  Add a new product to the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        item  body      AddItemRequest  true  "Cart item"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /carts/add [post]
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
    userID, _ := middleware.FromContext(r.Context())
    if userID == "" {
        http.Error(w, "unauthorized: userID missing in context", http.StatusUnauthorized)
        return
    }

    var req AddItemRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    item := &domain.CartItem{
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
    }

    if err := h.service.AddItem(userID, item); err != nil {
        if strings.Contains(err.Error(), "quantity") || strings.Contains(err.Error(), "stock") {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "item added"})
}


// @Summary      Remove Item from Cart
// @Description  Remove an item from the user's cart by product_id
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Param        product_id  query     string  true  "Product ID"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /carts/remove [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.FromContext(r.Context())
	if userID == "" {
		http.Error(w, "unauthorized: userID missing in context", http.StatusUnauthorized)
		return
	}
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "product_id required", http.StatusBadRequest)
		return
	}
	if err := h.service.RemoveItem(userID, productID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "item removed"})
}

// @Summary      Clear Cart
// @Description  Remove all items from the user's cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /carts/clear [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.FromContext(r.Context())
	if userID == "" {
		http.Error(w, "unauthorized: userID missing in context", http.StatusUnauthorized)
		return
	}
	if err := h.service.ClearCart(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "cart cleared"})
}
