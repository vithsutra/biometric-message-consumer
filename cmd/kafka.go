package main

import (
	"log"

	"github.com/IBM/sarama"
)

type kafkaConsumer struct {
	topic string
	group sarama.ConsumerGroup
}

func NewKafkaConsumer(brokers []string, groupId string, topic string) *kafkaConsumer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false

	group, err := sarama.NewConsumerGroup(brokers, groupId, config)

	if err != nil {
		log.Fatalln("error occurred while creating the consumer group, Error: ", err.Error())
	}

	return &kafkaConsumer{
		topic,
		group,
	}

}

func (c *kafkaConsumer) Close() {
	if c.group != nil {
		_ = c.group.Close()
	}
}
