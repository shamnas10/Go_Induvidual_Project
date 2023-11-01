package config

import (
	"datastream/logs"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

var Producer sarama.AsyncProducer

func InitKafkaProducer() {
	// Load environment variables from the .env file
	if err := godotenv.Load("/home/shamnas/Documents/test24/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Get the Kafka broker addresses from the environment variable
	brokers := os.Getenv("KAFKA_BROKERS")
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	var err error
	Producer, err = sarama.NewAsyncProducer([]string{brokers}, config)
	if err != nil {
		logs.Logger.Error("Error initializing Kafka producer", err)
	}

}

func CloseKafkaProducer() {
	if err := Producer.Close(); err != nil {
		logs.Logger.Error("Error closing Kafka producer", err)

	}
}

// Initialize a Kafka consumer
func InitializeKafkaConsumer(topic string) (sarama.Consumer, sarama.PartitionConsumer, error) {
	// Load environment variables from the .env file
	if err := godotenv.Load("/home/shamnas/Documents/test24/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	brokers := os.Getenv("KAFKA_BROKERS")
	// Create a Kafka consumer
	consumer, err := sarama.NewConsumer([]string{brokers}, consumerConfig)
	if err != nil {
		logs.Logger.Error("Error creating Kafka consumer", err)
		return nil, nil, err
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logs.Logger.Error("Error creating partition consumer", err)
		return nil, nil, err
	}

	return consumer, partitionConsumer, err
}
