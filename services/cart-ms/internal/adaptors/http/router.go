package http

import (
	//"ecom-api/pkg/middleware"
	"net/http"

	_ "github.com/babishagetaneh1992/ecom-api/services/cart-ms/docs"

	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Cart Microservice API
// @version         1.0
// @description     This is the Product service for the e-commerce system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8083
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// NewRouter sets up routes for order-ms
func NewRouter(handler *CartHandler) http.Handler {
	r := chi.NewRouter()

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/carts", func(r chi.Router) {
		r.With(middleware.AuthMiddleware).Get("/", handler.GetCart)
		r.With(middleware.AuthMiddleware).Post("/add", handler.AddItem)
		r.With(middleware.AuthMiddleware).Delete("/remove", handler.RemoveItem)
		r.With(middleware.AuthMiddleware).Delete("/clear", handler.ClearCart)
	})

	return r
}
