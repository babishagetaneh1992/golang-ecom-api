package main

import (
	"context"
	"syscall"

	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	//"user-microservice/adaptors/grpc/pb/user-microservice/services/user-ms/adaptors/grpc/pb"

	//"user-microservice/internal/adaptors/db"

	//"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/user-ms/services/user-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb"
	grpcAdapter "github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/adaptors/grpc"
	httpAdapter "github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/adaptors/http"

	//"user-microservice/internal/application"

	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
	//"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb/user-microservice/services/user-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/adaptors/db"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/application"
	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// @title           User Microservice API
// @version         1.0
// @description     This is the User service for the e-commerce system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.


func main() {
	// Load environment variables


	auth.InitJWT()

	// Get values from env
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	httpPort := os.Getenv("USER_HTTP_PORT")
	grpcPort := os.Getenv("USER_GRPC_PORT")

	if mongoURI == "" || dbName == "" {
		log.Fatal("Missing MONGO_URI or MONGO_DB_NAME in environment")
	}

	// MongoDB connection
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	dbConn := client.Database(dbName)

	// Layers
	repo := db.NewMongoUserRepository(dbConn)
	service := application.NewUserService(repo)

	// HTTP setup
	handler := httpAdapter.NewUserHandler(service)
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: httpAdapter.NewRouter(handler),
	}

	// gRPC setup
	grpcServer := grpc.NewServer()
	userGrpc := grpcAdapter.NewUserGrpcServer(service)
	pb.RegisterUserServiceServer(grpcServer, userGrpc)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	g := new(errgroup.Group) // run both http and grpc concurrently

	// HTTP server
	g.Go(func() error {
		fmt.Println("User-ms http server running on", httpPort)
		return httpServer.ListenAndServe()
	})

	// gRPC server
	g.Go(func() error {
		fmt.Println("User-ms grpc server running on", grpcPort)
		return grpcServer.Serve(lis)
	})

	// shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nshutting down server...") // syscall.SIGTERM

		// shutdown http server
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			log.Printf("HTTP server shutdown error: %v\n", err)
		}

		// stop grpc
		grpcServer.GracefulStop()
	}()

	// wait for either server to return an error
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
