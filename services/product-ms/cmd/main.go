package main

import (
	"context"
	"syscall"

	//"ecom-api/pkg/auth"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	//"product-microservice/adaptors/grpc/pb/product-microservice/services/product-ms/adaptors/grpc/pb"
	//"product-microservice/internal/adaptors/db"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/adaptors/grpc/pb"
	grpcAdapter "github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/adaptors/grpc"
	httpAdapter "github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/adaptors/http"

	//"product-microservice/internal/application"
	"time"

	_ "github.com/babishagetaneh1992/ecom-api/services/product-ms/docs"

	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
	//"github.com/babishagetaneh1992/ecom-api/services/product-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/product-ms/services/product-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/adaptors/db"
	"github.com/babishagetaneh1992/ecom-api/services/product-ms/internal/application"
	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// @title           Product Microservice API
// @version         1.0
// @description     This is the Product service for the e-commerce system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	
	
  auth.InitJWT()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	httpPort := os.Getenv("PRODUCT_HTTP_PORT")
	grpcPort := os.Getenv("PRODUCT_GRPC_ADDR")
	userMsAddr := os.Getenv("USER_GRPC_ADDR")

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
	repo := db.NewMongoProductRepository(dbConn)
	service := application.NewProductService(repo)

	// HTTP setup
	handler := httpAdapter.NewProductHandler(service)
	httpServer := http.Server{
		Addr:    httpPort,
		Handler: httpAdapter.NewRouter(handler),
	}

	// gRPC setup
	grpcServer := grpc.NewServer()
	productGrpc := grpcAdapter.NewProductGrpcServer(service)
	pb.RegisterProductServiceServer(grpcServer, productGrpc)

	// connect to user-ms
	conn, err := grpc.Dial(userMsAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect user-ms at %s: %v", userMsAddr, err)
	}
	defer conn.Close()

	// Example: you could initialize a user client here
	// userClient := productgrpc.NewUserClient(conn)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	g := new(errgroup.Group) // run both http and grpc concurrently

	// HTTP server
	g.Go(func() error {
		fmt.Println("Product-ms http server running on", httpPort)
		return httpServer.ListenAndServe()
	})

	// gRPC server
	g.Go(func() error {
		fmt.Println("Product-ms grpc server running on", grpcPort)
		return grpcServer.Serve(lis)
	})

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nshutting down server...")

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
