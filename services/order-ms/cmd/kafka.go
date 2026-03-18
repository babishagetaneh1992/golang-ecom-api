package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/adaptors/kafka"
	"github.com/babishagetaneh1992/ecom-api/services/order-ms/internals/ports"
)

// StartPaymentProcessedConsumer starts a background listener for payments.created events
func StartPaymentProcessedConsumer(ctx context.Context, kafkaConsumer *kafka.KafkaConsumer, service ports.OrderService) {
	fmt.Println("📡 Kafka Consumer listening for 'payments.created' events...")

	go func() {
		err := kafkaConsumer.ConsumeMessages("payments.created", func(msg *sarama.ConsumerMessage) error {
			// 1. Define the event structure (matching payment-ms)
			var event struct {
				PaymentID string  `json:"payment_id"`
				OrderID   string  `json:"order_id"`
				Amount    float64 `json:"amount"`
				Status    string  `json:"status"`
			}

			// 2. Decode the message
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("⚠️ Failed to parse Kafka message: %v\n", err)
				return nil
			}

			fmt.Printf("📥 Received Payment event: OrderID=%s, Status=%s\n", event.OrderID, event.Status)

			// 3. Update the order status in the database
			_, err := service.UpdateOrderStatus(ctx, event.OrderID, event.Status)
			if err != nil {
				fmt.Printf("❌ Failed to update order %s status: %v\n", event.OrderID, err)
			} else {
				fmt.Printf("✅ Order %s status updated to %s\n", event.OrderID, event.Status)
			}
			return nil
		})

		if err != nil {
			log.Printf("❌ Consumer loop error: %v\n", err)
		}
	}()
}
