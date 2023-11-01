package services

import (
	"datastream/config"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func SendToKafka(topic string, message []byte, Producer sarama.AsyncProducer) {
	config.Producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	go func() {
		for {
			select {
			case success := <-config.Producer.Successes():
				log.Printf("Produced message to topic %s, partition %d, offset %d\n,value %s",
					success.Topic, success.Partition,
					success.Offset, success.Value)
			case err := <-config.Producer.Errors():
				log.Printf("Failed to produce message: %v\n", err)
			}

		}
	}()
}

func ConsumeMessage(partitionConsumer sarama.PartitionConsumer) ([]byte, error) {
	select {
	case message := <-partitionConsumer.Messages():
		return message.Value, nil
	case err := <-partitionConsumer.Errors():
		return nil, err.Err
	case <-time.After(1 * time.Second):
		return nil, nil
	}
}
