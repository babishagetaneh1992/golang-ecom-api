package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/adaptors/kafka"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/domain"
	"github.com/babishagetaneh1992/ecom-api/services/payment-ms/internals/ports"
)

// StartOrderCreatedConsumer starts a background listener for order.created events
func StartOrderCreatedConsumer(ctx context.Context, kafkaConsumer *kafka.KafkaConsumer, service ports.PaymentService) {
	fmt.Println("📡 Kafka Consumer listening for 'order.created' events...")

	go func() {
		err := kafkaConsumer.ConsumeMessages("order.created", func(msg *sarama.ConsumerMessage) error {
			var event struct {
				OrderID string  `json:"order_id"`
				UserID  string  `json:"user_id"`
				Amount  float64 `json:"amount"`
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("⚠️ Failed to parse Kafka message: %v\n", err)
				return nil
			}

			fmt.Printf("📥 Received OrderCreated event: OrderID=%s, Amount=%.2f\n", event.OrderID, event.Amount)

			// Trigger payment processing automatically
			payment := &domain.Payment{
				OrderID: event.OrderID,
				UserID:  event.UserID,
				Amount:  event.Amount,
				Status:  "COMPLETED", // Simulated success
			}

			_, err := service.ProcessPayment(ctx, payment)
			if err != nil {
				fmt.Printf("❌ Failed to process payment for order %s: %v\n", event.OrderID, err)
			}
			return nil
		})

		if err != nil {
			log.Printf("❌ Consumer loop error: %v\n", err)
		}
	}()
}
