package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	//"ecom-api/pkg/auth"
	//"ecom-api/pkg/middleware"
	//"user-microservice/internal/domain"
	//"user-microservice/internal/ports"

	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/ports"
	"github.com/go-chi/chi/v5"
)

// Request structs for Swagger
type UserRegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

type UserLoginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(s ports.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      UserRegisterRequest  true  "User info"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	created, err := h.service.Register(&user)
	if err != nil {
		// duplicate email
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict) // 409 Conflict
			json.NewEncoder(w).Encode(map[string]string{
				"error": "email is already registered",
			})
			return
		}

		// validation errors
		if strings.Contains(err.Error(), "validation") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid user data",
			})
			return
		}

		// fallback
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "could not register user, please try again later",
		})
		return
	}

	resp := map[string]interface{}{
		"message": "user registered successfully",
		"user": map[string]interface{}{
			"id":    created.ID,
			"name":  created.Name,
			"email": created.Email,
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}


// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        login  body      UserLoginRequest  true  "Login info"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Router       /users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Requires JWT
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  domain.User
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.FromContext(r.Context())
	log.Println("Authenticated user:", userID)

	id := chi.URLParam(r, "id")

	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// ListUsers godoc
// @Summary      List all users
// @Description  Requires JWT
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   domain.User
// @Failure      401  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// ExistsUser godoc
// @Summary      Check if user exists
// @Description  Public endpoint
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]bool
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/exists [get]
func (h *UserHandler) ExistsUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	exists, err := h.service.Exists(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}
