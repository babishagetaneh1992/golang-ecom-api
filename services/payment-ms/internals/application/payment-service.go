package application

import (
	"context"
	"fmt"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/ports"
)

// PaymentServiceImplement implements ports.PaymentService
type PaymentServiceImplement struct {
	repo        ports.PaymentRepository
	orderClient ports.OrderClient // optional: if you want to call order-ms back
}

// constructor
func NewPaymentService(repo ports.PaymentRepository, orderClient ports.OrderClient) ports.PaymentService {
	return &PaymentServiceImplement{
		repo:        repo,
		orderClient: orderClient,
	}
}

// ProcessPayment simulates processing and persists result
func (s *PaymentServiceImplement) ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	// 1. Fetch order details to get total amount
	order, err := s.orderClient.GetOrder(ctx, payment.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order %s: %w", payment.OrderID, err)
	}

	// ✅ Always trust Order-MS total (prevents tampering)
	payment.Amount = order.Total
	payment.UserID = order.UserId

	// 2. Simulate payment gateway logic
	success := true // here you could integrate with Stripe/PayPal etc.
	if success {
		payment.Status = "COMPLETED"
	} else {
		payment.Status = "FAILED"
	}

	// 3. Persist payment record
	created, err := s.repo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to persist payment: %w", err)
	}

	// 4. Notify Order-MS of new status
	fmt.Println("🔎 Payment updating orderID=", payment.OrderID)
	fmt.Printf("📡 Sending status update to Order-MS: order=%s, newStatus=%s\n", payment.OrderID, payment.Status)

	if err := s.orderClient.UpdateOrderStatus(ctx, payment.OrderID, payment.Status); err != nil {
		fmt.Println("⚠️ warning: failed to notify order-ms:", err)
	} else {
		fmt.Printf("✅ Successfully updated order %s status in Order-MS to %s\n", payment.OrderID, payment.Status)
	}

	return created, nil
}


// InitPayment creates a PENDING payment for an order
func (s *PaymentServiceImplement) InitPayment(ctx context.Context, orderID string) (*domain.Payment, error) {
	payment := &domain.Payment{
		OrderID: orderID,
		Status:  "PENDING",
	}
	return s.repo.Create(ctx, payment)
}

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

