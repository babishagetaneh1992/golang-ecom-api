package grpc

import (

	//"cart-microservice/internal/domain"
	//"cart-microservice/internal/ports"
	"context"

	//"github.com/babishagetaneh1992/ecom-api/services/cart-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/cart-ms/services/cart-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/ports"
	//"product-microservice/adaptors/grpc/pb/product-microservice/services/product-ms/adaptors/grpc/pb"
)

type CartGrpcServer struct {
	pb.UnimplementedCartServiceServer
	service ports.CartService
}

func NewCartGrpcServer(s ports.CartService) *CartGrpcServer {
	return  &CartGrpcServer{service: s}
}

func (s *CartGrpcServer) AddItem(ctx context.Context, req *pb.AddItemRequest)(*pb.AddItemResponse, error) {
   item := domain.CartItem{
	ProductID: req.GetProductId(),
	Quantity: int(req.Quantity),
   }

   //call service to add item
   err := s.service.AddItem(req.GetUserId(), &item)
   if err != nil {
	return  nil, err
   }

   return  &pb.AddItemResponse{
	Message: "Item added to cart successfully",
   }, nil
}

func (s *CartGrpcServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	cart, err := s.service.GetCart(req.GetUserId())
	if err != nil {
		return nil, err
	}

	var pbItems []*pb.CartItem
	for _, item := range cart.Items {
		pbItems = append(pbItems, &pb.CartItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.GetCartResponse{
		Items: pbItems,
		Total: cart.Total,
	}, nil
}

func (s *CartGrpcServer) RemoveFromCart(ctx context.Context, req *pb.RemoveFromCartRequest) (*pb.RemoveFromCartResponse, error) {
	if err := s.service.RemoveItem(req.GetUserId(), req.GetProductId()); err != nil {
		return nil, err
	}

	return &pb.RemoveFromCartResponse{Message: "Item removed from cart successfully"}, nil
}

func (s *CartGrpcServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.ClearCartResponse, error) {
	if err := s.service.ClearCart(req.GetUserId()); err != nil {
		return nil, err
	}

	return &pb.ClearCartResponse{Message: "Cart cleared successfully"}, nil
}

