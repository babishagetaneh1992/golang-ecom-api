package application

import (
	"context"
	"fmt"

	//"order-microservice/internals/adaptors/grpc"
	//"order-microservice/internals/domain"
	//"order-microservice/internals/ports"

	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/adaptors/grpc"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/ports"
)

type OrderServiceImplement struct {
	repo          ports.OrderRepository
	cartClient    *grpc.CartClient
	paymentClient *grpc.PaymentClient
}

func NewOrderService(repo ports.OrderRepository, cartClient *grpc.CartClient, paymentClient *grpc.PaymentClient) ports.OrderService {
	return &OrderServiceImplement{
		repo:          repo,
		cartClient:    cartClient,
		paymentClient: paymentClient,
	}
}

// CreateOrderFromCart fetches the user's cart via gRPC and creates an order.
// It sets the order status = PENDING, saves it, then notifies Payment-MS.
// Payment-MS will later process the payment and call UpdateOrderStatus back.
func (s *OrderServiceImplement) CreateOrderFromCart(ctx context.Context, userID string) (*domain.Order, error) {
	// 1. Fetch cart
	cartResp, err := s.cartClient.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	if len(cartResp.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// 2. Build order
	var items []domain.OrderItem
	var totalPrice float64
	for _, ci := range cartResp.Items {
		price := ci.Price
		items = append(items, domain.OrderItem{
			ProductID: ci.ProductId,
			Quantity:  int(ci.Quantity),
			Price:     price,
		})
		totalPrice += float64(ci.Quantity) * price
	}

	order := &domain.Order{
		UserID:     userID,
		Items:      items,
		Total:    totalPrice,
		Status:     "PENDING",
	}

	// 3. Save order in DB
	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 4. Notify Payment-MS about new order
	// if msg, err := s.paymentClient.NotifyOrderCreated(ctx, createdOrder.ID); err != nil {
	// 	fmt.Println("warning: failed to notify payment-ms:", err)
	// } else {
	// 	fmt.Println("payment-ms response:", msg)
	// }

	paymentRes, err := s.paymentClient.NotifyOrderCreated(ctx, createdOrder.ID)
	if err != nil {
		fmt.Println("warning: failed to process payment:", err)
	} else {
		fmt.Println("payment-ms response:", paymentRes)
	}


	// 5. Clear cart
	if _, err := s.cartClient.ClearCart(ctx, userID); err != nil {
		fmt.Println("warning: failed to clear cart:", err)
	}

	return createdOrder, nil
}

func (s *OrderServiceImplement) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	order.Status = "PENDING"
	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// Notify Payment-MS about this order
	if msg, err := s.paymentClient.NotifyOrderCreated(ctx, createdOrder.ID); err != nil {
		fmt.Println("warning: failed to notify payment-ms:", err)
	} else {
		fmt.Println("payment-ms response:", msg)
	}

	return createdOrder, nil
}

func (s *OrderServiceImplement) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderServiceImplement) ListOrders(ctx context.Context) ([]*domain.Order, error) {
	return s.repo.List(ctx)
}

func (s *OrderServiceImplement) UpdateOrderStatus(ctx context.Context, id string, status string) (*domain.Order, error) {
	return s.repo.UpdateOrderStatus(ctx, id, status)
}

func (s *OrderServiceImplement) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
