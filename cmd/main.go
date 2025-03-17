package main

import "github.com/vithsutra/biometric-message-consumer/config"

func main() {

	config := config.InitConfig()

	db := NewDatabase(config.DatabaseUrl)

	db.CheckDatabaseConnection()

	defer db.CloseConnection()

	consumer := NewKafkaConsumer([]string{config.KafkaBroker1Address}, config.KafkaConsumerGroupId, config.KafkaTopic)

	defer consumer.Close()

	Start(db, consumer)

}
