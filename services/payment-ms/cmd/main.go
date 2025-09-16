package main

import (
	//"cart-microservice/internal/application"
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

	//"payment-microservice/adaptors/grpc/pb/payment-microservice/services/payment-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/adaptors/grpc/pb"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/adaptors/db"
	grpcAdapter "github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/adaptors/grpc"
	httpAdapter "github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/adaptors/http"

	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/application"

	_ "github.com/babishagetaneh1992/ecom-api/services/payment-ms/docs"

	//"github.com/babishagetaneh1992/ecom-api/services/payment-ms/adaptors/grpc/pb/github.com/babishagetaneh1992/ecom-api/services/payment-ms/services/payment-ms/adaptors/grpc/pb"
	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title           Payment Microservice API
// @version         1.0
// @description     Handles payments for orders in the e-commerce system.
// @host            localhost:8085
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {


 auth.InitJWT()
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	httpPort := os.Getenv("PAYMENT_HTTP_PORT")
	grpcPort := os.Getenv("PAYMENT_GRPC_ADDR")
	orderMSAddr := os.Getenv("ORDER_GRPC_ADDR")

	if mongoURI == "" || dbName == "" || grpcPort == "" {
		log.Fatal("❌ Missing required environment variables (MONGO_URI, MONGO_DB_NAME, PAYMENT_GRPC_PORT)")
	}

	// --- MongoDB ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)



	dbConn := client.Database(dbName)

	//grpc connection
	orderConn, err := grpc.Dial(orderMSAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to order-ms at %s: ", orderMSAddr)

	}
	defer orderConn.Close()
	orderClient := grpcAdapter.NewOrderClient(orderConn)



	repo := db.NewMongoPaymentRepository(dbConn)

	// --- Service ---
	service := application.NewPaymentService(repo, orderClient)

	//http set up
	handler := httpAdapter.NewPaymentHandler(service, orderClient)
	httpServer := &http.Server{
		Addr: httpPort,
		Handler: httpAdapter.NewPaymentRouter(handler),
	}

	// --- gRPC server ---
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, grpcAdapter.NewPaymentGrpcServer(service))

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	g := new(errgroup.Group)

	//http server
	g.Go(func() error {
		fmt.Println("✅ Payment http server is running on", httpPort)
		return httpServer.ListenAndServe()
	})

	// Start gRPC
	g.Go(func() error {
		fmt.Println("✅ Payment gRPC server running on", grpcPort)
		return grpcServer.Serve(lis)
	})

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\n🛑 Shutting down Payment service...")
		// shutdown http
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			log.Printf("Http shutdown error: %v\n", err)
		}
		grpcServer.GracefulStop()
		client.Disconnect(ctx)
	}()

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
