package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/babishagetaneh1992/ecom-api/services/notification-ms/internal/domain/service"
)

type Config struct {
	Broker     []string
	Topic      string
	SmsTopic   string
	EmailTopic string
	GroupID    string
}

type Consumer struct {
	group  sarama.ConsumerGroup
	config Config
}

func NewConsumer(config Config) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(config.Broker, saramaConfig)
	if err != nil {
		return nil, err
	}
	group, err := sarama.NewConsumerGroupFromClient(config.GroupID, client)
	if err != nil {
		return nil, err
	}
	return &Consumer{group: group, config: config}, nil
}

// Listen starts the consumer group loop
func (c *Consumer) Listen(ctx context.Context, topics []string) {
	handler := &ConsumerGroupHandler{}

	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.group.Consume(ctx, topics, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.group.Close()
}

// ConsumerGroupHandler represents a Sarama consumer group consumer
type ConsumerGroupHandler struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message recieved from topic: %s", message.Topic)

		switch message.Topic {
		case "order.created":
			var event service.OrderCreatedEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}
			
		log.Printf("Notification: New Order #%s placed by User #%s", event.OrderID, event.UserID)
		 
		// TODO: Add your notification logic here (e.g., sending SMS or Email)
		
		case "payment.created":
			var event service.PaymentCreatedEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}
			
		log.Printf("Notification: New Payment #%s for Order #%s", event.PaymentID, event.OrderID)
		 
		

		}
		     
		
		// TODO: Add your notification logic here (e.g., sending SMS or Email)
		
		session.MarkMessage(message, "")
	}

	return nil
}




