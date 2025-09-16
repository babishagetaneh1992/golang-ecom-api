package grpc

import (
	"context"
	//"product-microservice/adaptors/grpc/pb/user-microservice/services/product-ms/adaptors/grpc/pb"
	//"product-microservice/adaptors/grpc/pb/product-microservice/services/product-ms/adaptors/grpc/pb"
	//"product-microservice/internal/domain"
	//"product-microservice/internal/ports"

	//"github.com/babishagetaneh1992/ecom-api/services/product-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/product-ms/services/product-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/ports"
)

type ProductGrpcServer struct {
	pb.UnimplementedProductServiceServer
	service ports.ProductService
}

func NewProductGrpcServer(s ports.ProductService) *ProductGrpcServer {
	return  &ProductGrpcServer{service: s}
}

func (s *ProductGrpcServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := &domain.Product{
		Name: req.GetProduct().GetName(),
		Description: req.GetProduct().GetDescription(),
		Price: req.GetProduct().GetPrice(),
		Stock: int(req.GetProduct().GetStock()),
	}

	p, err := s.service.CreateNewProduct(ctx, product)
	if err != nil {
		return  nil, err
	}

	 return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id: p.ID,
			Name: p.Name,
			Description: p.Description,
			Price: p.Price,
			Stock: int32(p.Stock),
		},
	 }, nil
}


func (s *ProductGrpcServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		return  nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: product.Price,
			Stock: int32(product.Stock),
		},
	}, nil
}

func (s *ProductGrpcServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.service.ListProducts(ctx)
	if err != nil {
		return  nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id: p.ID,
			Name: p.Name,
			Description:  p.Description,
			Price: p.Price,
			Stock: int32(p.Stock),
		})
	}

	return  &pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}


func (s *ProductGrpcServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	product := &domain.Product{
		ID:          req.Product.Id,
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       req.Product.Price,
		Stock:       int(req.Product.Stock),
	}

	p, err := s.service.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
		},
	}, nil
}

func (s *ProductGrpcServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := s.service.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}