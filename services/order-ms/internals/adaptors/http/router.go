package http

import (
	"net/http"
	//"ecom-api/pkg/middleware"
	"github.com/babishagetaneh1992/ecom-api/pkg/middleware"
	userPb "github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter sets up the routes for the Order microservice
// @title           Order Microservice API
// @version         1.0
// @description     This service manages orders in the e-commerce system.
// @host            localhost:8083
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func NewRouter(handler *OrderHandler, userClient userPb.UserServiceClient) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.GRPCAuthMiddleware(userClient))

		r.With(chiMiddleware.AllowContentType("application/json")).Post("/",handler.CreateOrder)
		r.With(chiMiddleware.AllowContentType("application/json")).Put("/{id}",handler.UpdateOrderStatus)
		r.Get("/", handler.ListOrders)
		r.Get("/{id}", handler.GetOrder)
		r.Delete("/{id}", handler.DeleteOrder)

	})

	return  r
}