package grpc

import (
	"context"
	"fmt"

	//"order-microservice/adaptors/grpc/pb/order-microservice/services/order-ms/adaptors/grpc/pb"
	//"payment-microservice/internals/domain"

	//"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/order-ms/services/order-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client pb.OrderServiceClient
}

func NewOrderClient(conn *grpc.ClientConn) *OrderClient {
	return &OrderClient{
		client: pb.NewOrderServiceClient(conn),
	}
}

// UpdateOrderStatus calls Order-MS to update order status
func (c *OrderClient) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	_, err := c.client.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		Id:     orderID,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}


func (c *OrderClient) GetOrder(ctx context.Context, orderID string) (*pb.Order, error) {
	resp, err := c.client.GetOrder(ctx, &pb.GetOrderRequest{Id: orderID})
	if err != nil {
		return nil, err
	}

	return resp.Order, nil
}