package ports

import (
	"context"

	"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb"
	//"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/order-ms/services/order-ms/adaptors/grpc/pb"
	//"order-microservice/adaptors/grpc/pb/order-microservice/services/order-ms/adaptors/grpc/pb"
	//"payment-microservice/internals/domain"
)

type OrderClient interface {
    UpdateOrderStatus(ctx context.Context, orderID, status string) error
	GetOrder(ctx context.Context, orderID string) (*pb.Order, error)
}