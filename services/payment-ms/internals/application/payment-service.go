package application

import (
	"context"
	"fmt"

	//"github.com/IBM/sarama"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/adaptors/kafka"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/ports"
)

// PaymentServiceImplement implements ports.PaymentService
type PaymentServiceImplement struct {
	repo        ports.PaymentRepository
	//orderClient ports.OrderClient // optional: if you want to call order-ms back
	kafkaProducer kafka.KafkaProducer
}

// constructor
func NewPaymentService(repo ports.PaymentRepository, kafkaProducer kafka.KafkaProducer) ports.PaymentService {
	return &PaymentServiceImplement{
		repo:        repo,
		//orderClient: orderClient,
		kafkaProducer: kafkaProducer,
	}
}

func (s *PaymentServiceImplement) ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	// 1. Persist payment record first to get the ID
	created, err := s.repo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to persist payment: %w", err)
	}

	// 2. Now produce the event using the created payment's information
	event := &kafka.PaymentCreatedEvent{
		PaymentID: created.ID,
		OrderID:   created.OrderID,
		Amount:    created.Amount,
		Status:    created.Status,
	}

	if err := s.kafkaProducer.ProducePaymentCreatedEvent(event); err != nil {
		return nil, fmt.Errorf("failed to produce payment created event: %w", err)
	}

	return created, nil
}


// InitPayment creates a PENDING payment for an order
// func (s *PaymentServiceImplement) InitPayment(ctx context.Context, orderID string) (*domain.Payment, error) {
// 	payment := &domain.Payment{
// 		OrderID: orderID,
// 		Status:  "COMPLETED",
// 	}
// 	return s.repo.Create(ctx, payment)
// }

// GetPayment fetches a payment by ID
func (s *PaymentServiceImplement) GetPayment(ctx context.Context, id string) (*domain.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

// ListPayments returns all payments
func (s *PaymentServiceImplement) ListPayments(ctx context.Context) ([]*domain.Payment, error) {
	return s.repo.List(ctx)
}

// UpdatePaymentStatus changes the status of a payment
func (s *PaymentServiceImplement) UpdatePaymentStatus(ctx context.Context, id string, status string) (*domain.Payment, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

// DeletePayment removes a payment record
func (s *PaymentServiceImplement) DeletePayment(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// NotifyOrderCreated is called by Order-MS when a new order is placed.
// It creates a PENDING payment record but does not process it yet.
func (s *PaymentServiceImplement) NotifyOrderCreated(ctx context.Context, orderID string) error {
    fmt.Printf("Order %s created, payment initialization deferred\n", orderID)
    return nil
}

