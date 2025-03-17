package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Variables struct {
	DatabaseUrl          string
	KafkaBroker1Address  string
	KafkaConsumerGroupId string
	KafkaTopic           string
}

func InitConfig() *Variables {
	serverMode := os.Getenv("SERVER_MODE")

	if serverMode != "dev" && serverMode != "prod" {
		log.Fatalln("please set SERVER_MODE to dev or prod")
	}

	if serverMode == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("failed to load the .env file, Error: ", err.Error())
		}
	}

	variable := new(Variables)

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatalln("missing or empty DATABASE_URL env variable")
	}

	kafkaBroker1Address := os.Getenv("KAFKA_BROKER_1_ADDRESS")

	if kafkaBroker1Address == "" {
		log.Fatalln("missing or empty KAFKA_BROKER_1_ADDRESS env variable")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC_NAME")

	if kafkaTopic == "" {
		log.Fatalln("missing KAFKA_TOPIC_NAME env variable")
	}

	kafkaConsumerGroupId := os.Getenv("KAFKA_CONSUMER_GROUP_ID")

	if kafkaTopic == "" {
		log.Fatalln("missing KAFKA_CONSUMER_GROUP_ID env variable")
	}

	variable.DatabaseUrl = dbUrl
	variable.KafkaBroker1Address = kafkaBroker1Address
	variable.KafkaConsumerGroupId = kafkaConsumerGroupId
	variable.KafkaTopic = kafkaTopic

	return variable
}
