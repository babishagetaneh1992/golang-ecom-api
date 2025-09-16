package http

import (
	"net/http"
	appMiddleware "github.com/babishagetaneh1992/ecom-api/pkg/middleware"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/babishagetaneh1992/ecom-api/services/product-ms/docs" // Swagger docs
)

// @title           Product Microservice API
// @version         1.0
// @description     This is the Product service for the e-commerce system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func NewRouter(handler *ProductHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Product routes
	r.Route("/products", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)

		r.Post("/", handler.CreateProduct)
		r.Get("/", handler.ListProducts)
		r.Get("/{id}", handler.GetProduct)
		r.Put("/{id}", handler.UpdateProduct)
		r.Delete("/{id}", handler.DeleteProduct)
	})

	return r
}
