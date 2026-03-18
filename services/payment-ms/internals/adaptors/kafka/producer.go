package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)


type PaymentCreatedEvent struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Amount    float64 `json:"amount"`
	Status    string `json:"status"`
}


type KafkaProducer struct {
	client sarama.AsyncProducer
	topic string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			fmt.Println("Error producing message:", err)
		}
	}()

	go func() {
		for msg := range producer.Successes() {
			fmt.Println("Message produced successfully:", msg)
		}
	}()

	return &KafkaProducer{
		client: producer,
		topic:  topic,
	}, nil
}



func (p *KafkaProducer) ProducePaymentCreatedEvent(event *PaymentCreatedEvent) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.ProduceMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	})
}


func (p *KafkaProducer) ProduceMessage(msg *sarama.ProducerMessage) error {
	p.client.Input() <- msg
	return nil
}


func (p *KafkaProducer) Close() error {
	return p.client.Close()
}