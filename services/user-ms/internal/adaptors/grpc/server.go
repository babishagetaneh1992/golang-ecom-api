package grpc

import (
	"context"

	//"user-microservice/internal/ports"

	//"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb/user-microservice/services/user-ms/adaptors/grpc/pb"
	//"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/user-ms/services/user-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/ports"
)

type UserGrpcServer struct {
	pb.UnimplementedUserServiceServer
	service ports.UserService
}


func NewUserGrpcServer(service ports.UserService) *UserGrpcServer {
	return  &UserGrpcServer{service: service}
}


func (s *UserGrpcServer) RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &domain.User{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	}

	createdUser, err := s.service.Register(user)
	if err != nil {
		return nil, err
	}

	return  &pb.CreateUserResponse{
		User: &pb.User{
			Id: createdUser.ID,
			Name: createdUser.Name,
			Email: createdUser.Email,
			Password: createdUser.Password,
			CreatedAt: createdUser.CreateAt.String(),
		},
	}, nil
}


func (s *UserGrpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.service.GetUser(req.Id)
	if err != nil {
		return  nil, err
	}

	return  &pb.GetUserResponse{
		User: &pb.User{
			Id: user.ID,
			Name: user.Name,
			Email: user.Email,
			Password: user.Password,
			CreatedAt: user.CreateAt.String(),
		},
	}, nil
}

func (s *UserGrpcServer) ListUsers(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	users, err := s.service.ListUsers()
	if err != nil {
		return  nil, err
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id: u.ID,
			Name: u.Name,
			Email: u.Email,
			Password: u.Password,
			CreatedAt: u.CreateAt.String(),
		})
	}

	return  &pb.ListUserResponse{
		Users: pbUsers,

	}, nil
}



func (s *UserGrpcServer) ExistUser(ctx context.Context, req *pb.ExistRequest)(*pb.ExistResponse, error) {
	exists, err := s.service.Exists(req.Id)
	if err != nil {
		return  nil, err
	}

	return  &pb.ExistResponse{
		Exists: exists,
	}, nil
}
