package application

import (
	//"cart-microservice/internal/adaptors/grpc"
	//"cart-microservice/internal/domain"
	//"cart-microservice/internal/ports"
	"fmt"

	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/adaptors/grpc"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/ports"
)

type CartServiceImplement struct {
	repo ports.CartRepository
	productClient  *grpc.ProdctClient
	
}

func NewCartService(r ports.CartRepository, productClient *grpc.ProdctClient) ports.CartService {
	return  &CartServiceImplement{
		repo: r,
		productClient: productClient ,
	}
}


func (s *CartServiceImplement) AddItem(userID string, item *domain.CartItem) error {
	// 1. Validate quantity
	if item.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	// 2. Call Product-MS via gRPC
	product, err := s.productClient.GetProduct(item.ProductID)
	if err != nil {
		return fmt.Errorf("failed to fetch product %s: %w", item.ProductID, err)
	}

	// 3. Check stock
	if int32(item.Quantity) > product.Product.Stock {
		return fmt.Errorf("requested quantity %d exceeds available stock %d",
			item.Quantity, product.Product.Stock)
	}

	// 4. Copy price & name from product
	item.Price = product.Product.Price
	item.Name = product.Product.Name

	// 5. Save to repo
	if err := s.repo.AddItem(userID, *item); err != nil {
		return fmt.Errorf("failed to save cart item: %w", err)
	}

	return nil
}


func (s *CartServiceImplement) GetCart(userID string) (*domain.Cart, error) {
	cart, err :=  s.repo.GetCart(userID)
	if err != nil {
		return  nil, fmt.Errorf("failed to get cart: %w", err)
	}

	var total float64
	for _, item := range cart.Items {
		total += float64(item.Quantity) * item.Price
	}

	cart.Total = total

	return  cart, nil
}

func (s *CartServiceImplement) RemoveItem(userID, productID string) error {
	if err := s.repo.RemoveItem(userID, productID); err != nil {
		return  fmt.Errorf("failed to remove item: %w", err)
	}
    
	return  nil
}

func (s *CartServiceImplement) ClearCart(userID string) error {
	if err := s.repo.ClearCart(userID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}