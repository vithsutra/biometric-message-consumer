package main

import (
	"context"
	"log"

	"github.com/vithsutra/biometric-message-consumer/handler"
	"github.com/vithsutra/biometric-message-consumer/repository"
)

func Start(database *database, kafkaConsumer *kafkaConsumer) {
	postgresRepo := repository.NewPostgresRepository(database.conn)

	consumerHandler := handler.NewConsumerGroupHandler(postgresRepo)

	log.Println("consumer running..")
	for {
		if err := kafkaConsumer.group.Consume(context.Background(), []string{kafkaConsumer.topic}, consumerHandler); err != nil {
			log.Println("error occurred while conuming the message, Error: ", err.Error())
			continue
		}
	}
}
