package grpc

import (
	"context"
	//"order-microservice/adaptors/grpc/pb/order-microservice/services/order-ms/adaptors/grpc/pb"
	//"order-microservice/internals/domain"
	//"order-microservice/internals/ports"

	//"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/order-ms/services/order-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/ports"
)

type OrderGrpcServer struct {
	pb.UnimplementedOrderServiceServer
	service ports.OrderService
}

func NewOrderGrpcServer(s *ports.OrderService) *OrderGrpcServer {
    return  &OrderGrpcServer{service: *s}
}

func (s *OrderGrpcServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		}
	}

	order := &domain.Order{
		UserID: req.UserId,
		Items:  items,
		Status: "PENDING",
	}

	created, err := s.service.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Order: toProto(created),
	}, nil
}

func (s *OrderGrpcServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{Order: toProto(order)}, nil
}

func (s *OrderGrpcServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.service.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListOrdersResponse{}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return resp, nil
}

func (s *OrderGrpcServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	// Call service to update status
	_, err := s.service.UpdateOrderStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}

	// After updating, fetch the updated order to return
	order, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrderStatusResponse{
		Order: toProto(order),
	}, nil
}


func (s *OrderGrpcServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := s.service.DeleteOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{Message: "order deleted successfully"}, nil
}

// helper to convert domain → proto
func toProto(order *domain.Order) *pb.Order {
	if order == nil {
		return nil
	}

	items := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		}
	}

	return &pb.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Items:  items,
		Total: order.Total,
		Status: order.Status,
	}
}
