package application

import (
	"context"
	"fmt"

	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/ports"
)

type ProductServiceimplement struct {
	repo ports.ProductRepository
}

func NewProductService(r ports.ProductRepository) ports.ProductService {
  return  &ProductServiceimplement{repo: r}
}

func (s *ProductServiceimplement) CreateNewProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	// Validation
	if p.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if p.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}
	if p.Stock < 0 {
		return nil, fmt.Errorf("stock cannot be negative")
	}

	// If validation passes → Save to DB
	return s.repo.CreateProduct(ctx, p)
}


func (s *ProductServiceimplement) GetProduct(ctx context.Context, id string)(*domain.Product, error) {
	return s.repo.FindByID(ctx,id)
}

func (s *ProductServiceimplement) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return  s.repo.FindAll(ctx)
}

func (s *ProductServiceimplement) UpdateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	return  s.repo.Update(ctx, p)
}

func (s *ProductServiceimplement) DeleteProduct(ctx context.Context, id string) error {
	return  s.repo.Delete(ctx, id)
}
