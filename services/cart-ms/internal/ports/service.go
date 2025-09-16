package ports

import "github.com/babishagetaneh1992/ecom-api/services/cart-ms/internal/domain"

//import "cart-microservice/internal/domain"

type CartService interface {
	GetCart(userID string) (*domain.Cart, error)
	AddItem(userID string, item *domain.CartItem) error 
	RemoveItem(userID, productID string) error
	ClearCart(userID string) error
}