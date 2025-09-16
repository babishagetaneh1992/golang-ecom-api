package http

import (
	"encoding/json"
	"net/http"

	//"product-microservice/internal/domain"
	//"product-microservice/internal/ports"

	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/ports"
	"github.com/go-chi/chi"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	service ports.ProductService
}

func NewProductHandler(s ports.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// Request DTOs for Swagger
type ProductCreateRequest struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"High-end gaming laptop"`
	Price       float64 `json:"price" example:"1299.99"`
	Stock       int     `json:"stock" example:"10"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"Updated description"`
	Price       float64 `json:"price" example:"1199.99"`
	Stock       int     `json:"stock" example:"15"`
}

// @Summary      Create product
// @Description  Add a new product (requires JWT)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product  body      ProductCreateRequest  true  "Product info"
// @Success      201  {object}  domain.Product
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	// Map request → domain model
	p := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	product, err := h.service.CreateNewProduct(r.Context(), &p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}


	response := map[string]interface{}{
		"message": "product created successfully",
		"product": map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Get product by ID
// @Description  Fetch a product by ID (requires JWT)
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      404  {object}  map[string]string
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// @Summary      List products
// @Description  Get all products (requires JWT)
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.Product
// @Failure      500  {object}  map[string]string
// @Router       /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// @Summary      Update product
// @Description  Update product by ID (requires JWT)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string              true  "Product ID"
// @Param        product  body      ProductUpdateRequest true "Updated product"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := domain.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	product, err := h.service.UpdateProduct(r.Context(), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// @Summary      Delete product
// @Description  Delete product by ID (requires JWT)
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Product deleted successfully"}`))
}
