package grpc

import (
	"context"
	//"product-microservice/adaptors/grpc/pb/product-microservice/services/product-ms/adaptors/grpc/pb"

	//"github.com/babishagetaneh1992/ecom-api/services/product-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/product-ms/services/product-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/adaptors/grpc/pb"
	"google.golang.org/grpc"
)



type ProdctClient struct {
	client pb.ProductServiceClient
}

func NewProductClient(conn *grpc.ClientConn ) *ProdctClient {
  return  &ProdctClient{
	client: pb.NewProductServiceClient(conn),
  }
}

func (c *ProdctClient) GetProduct(id string) (*pb.GetProductResponse, error) {
	res, err := c.client.GetProduct(context.Background(), &pb.GetProductRequest{Id: id})
	if err != nil {
		return  nil, err
	}

	return res, nil
}