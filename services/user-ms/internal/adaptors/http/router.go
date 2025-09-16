package http

import (
	"net/http"

	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	//"ecom-api/pkg/middleware"
	_ "github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/adaptors/http/docs" // Swagger docs
)

// NewRouter configures and returns a Chi router with all user routes and Swagger
func NewRouter(handler *UserHandler) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chi_middleware.Logger)    // logs requests
	r.Use(chi_middleware.Recoverer) // recovers from panics
	r.Use(chi_middleware.AllowContentType("application/json"))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Swagger route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// User routes
	r.Route("/users", func(r chi.Router) {

		// ✅ Public routes (no auth)
		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.Login)
		r.Get("/{id}/exists", handler.ExistsUser)

		// ✅ Protected routes (require JWT)
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthMiddleware)

			protected.Get("/", handler.ListUsers)
			protected.Get("/{id}", handler.GetUser)
		})
	})

	return r
}
