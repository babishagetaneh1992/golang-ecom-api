package grpc

import (
	//"cart-microservice/adaptors/grpc/pb/cart-microservice/services/cart-ms/adaptors/grpc/pb"
	"context"

	//"github.com/babishagetaneh1992/ecom-api/services/cart-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/cart-ms/services/cart-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/adaptors/grpc/pb"
	"google.golang.org/grpc"
)

type CartClient struct {
	client pb.CartServiceClient
}

func NewCartClient(conn *grpc.ClientConn) *CartClient {
	return  &CartClient{
		client: pb.NewCartServiceClient(conn),
	}
}

func (c *CartClient) GetCart(ctx context.Context, userID string) (*pb.GetCartResponse, error) {
	res, err := c.client.GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	return  res, nil
}


func (c *CartClient) ClearCart(ctx context.Context, userID string) (*pb.ClearCartResponse, error) {
	return c.client.ClearCart(ctx, &pb.ClearCartRequest{UserId: userID})
}


