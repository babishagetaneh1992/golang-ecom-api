package ports

import (
	"context"

	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	//"payment-microservice/internals/domain"
)

type PaymentService interface {
    ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
    InitPayment(ctx context.Context, orderID string) (*domain.Payment, error)
    GetPayment(ctx context.Context, id string) (*domain.Payment, error)
    ListPayments(ctx context.Context) ([]*domain.Payment, error)
    UpdatePaymentStatus(ctx context.Context, id string, status string) (*domain.Payment, error)
    DeletePayment(ctx context.Context, id string) error
	NotifyOrderCreated(ctx context.Context, orderID string) error
}
