package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/babishagetaneh1992/ecom-api/services/notification-ms/internal/adapters/kafka"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("No .env file found, relying on system env variables")
	}

	brokers := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("NOTIFICATION_KAFKA_GROUP_ID")

	if brokers == "" || groupID == "" {
		log.Fatal("❌ Missing required Kafka configuration (KAFKA_BROKERS or NOTIFICATION_KAFKA_GROUP_ID)")
	}

	// Initialize Kafka Consumer
	config := kafka.Config{
		Broker:  []string{brokers},
		GroupID: groupID,
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Topics to listen to
	topics := []string{"order.created", "payments.created"}

	fmt.Println("🚀 Notification Service is starting...")
	fmt.Printf("📡 Listening to topics: %v\n", topics)
	fmt.Printf("👥 Consumer Group: %s\n", groupID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start listening in the background
	consumer.Listen(ctx, topics)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\n🛑 Shutting down Notification service...")
	cancel() // Signal the consumer to stop
}
