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

func NewPaymentRouter(handler *PaymentHandler, userClient userPb.UserServiceClient) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Swagger route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/payments", func(r chi.Router) {
		r.Use(middleware.GRPCAuthMiddleware(userClient))

		// No request body, only orderID in path + Bearer token
		r.Post("/{order_id}", handler.CreatePayment)
	})

	return r
}
