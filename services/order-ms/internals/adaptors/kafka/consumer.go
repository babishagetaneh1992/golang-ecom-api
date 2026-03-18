package kafka

import "github.com/IBM/sarama"

type KafkaConsumer struct {
	client *sarama.Consumer
}

func NewKafkaConsumer(brokers []string, topic string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		client: &consumer,
	}, nil
}



func (c *KafkaConsumer) ConsumeMessages(topic string, handler func(*sarama.ConsumerMessage) error) error {
	partitionConsumer, err := (*c.client).ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		if err := handler(msg); err != nil {
			return err
		}
	}
	return nil
}


func (c *KafkaConsumer) Close() {
	(*c.client).Close()
}