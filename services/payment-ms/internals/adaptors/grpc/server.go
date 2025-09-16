package grpc

import (
	"context"
	"fmt"

	//"payment-microservice/adaptors/grpc/pb/payment-microservice/services/payment-ms/adaptors/grpc/pb"
	//"github.com/babishagetaneh1992/ecom-api/services/payment-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/payment-ms/services/payment-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/ports"
)

type PaymentGrpcServer struct {
	pb.UnimplementedPaymentServiceServer
	service ports.PaymentService
}

func NewPaymentGrpcServer(service ports.PaymentService) *PaymentGrpcServer{
	return  &PaymentGrpcServer{service: service}
}

func (s *PaymentGrpcServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	payment := &domain.Payment{
		OrderID: req.GetOrderId(),
		UserID:  req.GetUserId(),
		Amount:  req.GetAmount(),
		Status:  "PENDING",
		//Method:  req.GetMethod(),
	}

	createdPayment, err := s.service.ProcessPayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	return &pb.ProcessPaymentResponse{
		Payment: &pb.Payment{
			Id:      createdPayment.ID,
			OrderId: createdPayment.OrderID,
			UserId:  createdPayment.UserID,
			Amount:  createdPayment.Amount,
			Status:  createdPayment.Status,
			//Method:  createdPayment.Method,
		},
	}, nil
}

func (s *PaymentGrpcServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	payment, err := s.service.GetPayment(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &pb.GetPaymentResponse{
		Payment: &pb.Payment{
			Id:      payment.ID,
			OrderId: payment.OrderID,
			UserId:  payment.UserID,
			Amount:  payment.Amount,
			Status:  payment.Status,
			//Method:  payment.Method,
		},
	}, nil
}

func (s *PaymentGrpcServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	payments, err := s.service.ListPayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	var pbPayments []*pb.Payment
	for _, p := range payments {
		pbPayments = append(pbPayments, &pb.Payment{
			Id:      p.ID,
			OrderId: p.OrderID,
			UserId:  p.UserID,
			Amount:  p.Amount,
			Status:  p.Status,
			//Method:  p.Method,
		})
	}

	return &pb.ListPaymentsResponse{Payments: pbPayments}, nil
}

func (s *PaymentGrpcServer) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error) {
	updatedPayment, err := s.service.UpdatePaymentStatus(ctx, req.GetId(), req.GetStatus())
	if err != nil {
		return nil, fmt.Errorf("failed to update payment with id %s: %w", req.GetId(), err)
	}

	return &pb.UpdatePaymentStatusResponse{
		Payment: &pb.Payment{
			Id:      updatedPayment.ID,
			OrderId: updatedPayment.OrderID,
			UserId:  updatedPayment.UserID,
			Amount:  updatedPayment.Amount,
			Status:  updatedPayment.Status,
			//Method:  updatedPayment.Method,
		},
	}, nil
}




// DeletePayment deletes a payment by ID
func (s *PaymentGrpcServer) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*pb.DeletePaymentResponse, error) {
	if err := s.service.DeletePayment(ctx, req.GetId()); err != nil {
		return nil, fmt.Errorf("failed to delete payment: %w", err)
	}

	return &pb.DeletePaymentResponse{Message: "Payment deleted successfully"}, nil
}

// adaptors/grpc/server.go (payment-ms)

func (s *PaymentGrpcServer) NotifyOrderCreated(ctx context.Context, req *pb.NotifyOrderRequest) (*pb.NotifyOrderResponse, error) {
    orderID := req.GetOrderId()
    if orderID == "" {
        return nil, fmt.Errorf("order_id is required")
    }

   

    // ✅ NEW: delegate to application-layer NotifyOrderCreated
    if err := s.service.NotifyOrderCreated(ctx, orderID); err != nil {
        return nil, fmt.Errorf("failed to handle NotifyOrderCreated for order %s: %w", orderID, err)
    }

    return &pb.NotifyOrderResponse{
        Message: fmt.Sprintf("Payment record initialized for order %s", orderID),
    }, nil
}
