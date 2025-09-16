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
	"time"

	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
	//"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/order-ms/services/order-ms/adaptors/grpc/pb"

	//"order-microservice/adaptors/grpc/pb/order-microservice/services/order-ms/adaptors/grpc/pb"
	//"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/adaptors/db"
	grpcAdapter "github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/adaptors/grpc"
	httpAdapter "github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/adaptors/http"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/application"

	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/babishagetaneh1992/ecom-api/services/order-ms/docs"
)

// @title           Order Microservice API
// @version         1.0
// @description     This is the Order service for the e-commerce system.
// @host            localhost:8084
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.


func main() {

   auth.InitJWT()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	httpPort := os.Getenv("ORDER_HTTP_PORT")
	grpcPort := os.Getenv("ORDER_GRPC_ADDR")
	cartMsAddr := os.Getenv("CART_GRPC_ADDR")
	paymentMsAddr := os.Getenv("PAYMENT_GRPC_ADDR")

	if mongoURI == "" || dbName == "" || httpPort == "" || grpcPort == "" {
		log.Fatal("❌ Missing required env vars: MONGO_URI, MONGO_DB_NAME, ORDER_HTTP_PORT, ORDER_GRPC_PORT")
	}

	// --- MongoDB ---
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

	// --- gRPC clients for dependencies ---
	// Cart-MS
	cartConn, err := grpc.Dial(cartMsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to cart-ms at %s: %v", cartMsAddr, err)
	}
	defer cartConn.Close()
	cartClient := grpcAdapter.NewCartClient(cartConn)

	// Payment-MS
	paymentConn, err := grpc.Dial(paymentMsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment-ms at %s: %v", paymentMsAddr, err)
	}
	defer paymentConn.Close()
	paymentClient := grpcAdapter.NewPaymentClient(paymentConn)

	// --- Service ---
	repo := db.NewMongoOrderRepository(dbConn)
	service := application.NewOrderService(repo, cartClient, paymentClient)

	// --- HTTP setup ---
	handler := httpAdapter.NewOrderHandler(service)
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: httpAdapter.NewRouter(handler),
	}

	// --- gRPC setup ---
	grpcServer := grpc.NewServer()
	orderGrpc := grpcAdapter.NewOrderGrpcServer(&service)
	pb.RegisterOrderServiceServer(grpcServer, orderGrpc)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	g := new(errgroup.Group)

	// HTTP server
	g.Go(func() error {
		fmt.Println("✅ Order HTTP server running on", httpPort)
		return httpServer.ListenAndServe()
	})

	// gRPC server
	g.Go(func() error {
		fmt.Println("✅ Order gRPC server running on", grpcPort)
		return grpcServer.Serve(lis)
	})

	// --- Graceful shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\n🛑 Shutting down Order service...")

		// shutdown HTTP
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			log.Printf("HTTP shutdown error: %v\n", err)
		}

		// shutdown gRPC
		grpcServer.GracefulStop()
	}()

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
