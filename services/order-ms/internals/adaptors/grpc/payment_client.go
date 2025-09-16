package grpc

import (
	"context"
	"log"

	//"payment-microservice/adaptors/grpc/pb/payment-microservice/services/payment-ms/adaptors/grpc/pb"

	//"order-microservice/services/order-ms/adaptors/grpc/pb/payment" // adjust import path

	//"github.com/babishagetaneh1992/ecom-api/services/payment-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/payment-ms/services/payment-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/adaptors/grpc/pb"
	"google.golang.org/grpc"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *PaymentClient {
	return &PaymentClient{
		client: pb.NewPaymentServiceClient(conn),
	}
}

func (p *PaymentClient) ProcessPayment(ctx context.Context, orderID, userID string, amount float64) (*pb.Payment, error) {
	resp, err := p.client.ProcessPayment(ctx, &pb.ProcessPaymentRequest{
		OrderId: orderID,
		UserId:  userID,
		Amount:  amount,
	})
	if err != nil {
		log.Printf("failed to process payment: %v", err)
		return nil, err
	}
	return resp.Payment, nil
}

func (p *PaymentClient) NotifyOrderCreated(ctx context.Context, orderID string) (string, error) {
	resp, err := p.client.NotifyOrderCreated(ctx, &pb.NotifyOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		return "", err
	}
	return resp.GetMessage(), nil
}
