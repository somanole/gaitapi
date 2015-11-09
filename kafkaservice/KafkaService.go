// KafkaService
package kafkaservice

import (
	"github.com/Shopify/sarama"
	"github.com/somanole/gaitapi/constants"
	"log"
	"fmt"
)

func ProduceDummyMessage() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{constants.KAFKA_MASTER}, config)
	if err != nil {
	    log.Printf(fmt.Sprintf("KafkaService.ProduceDummyMessage() - Error: %v", err.Error()))
	}
	
	message := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder("ecca3bd1-75b0-4181-99c9-96817cd5753e,1445808714,1445809489")}
	producer.Input() <- message
	producer.AsyncClose() 
	
	log.Printf("KafkaService.ProduceDummyMessage() - Successfully produced dummy message!")
	
	return err
}

func ProduceMessage(receivedMessage string) error {
	log.Printf(fmt.Sprintf("KafkaService.ProduceMessage() - received message: %v", receivedMessage))
	
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{constants.KAFKA_MASTER}, config)
	if err != nil {
	    log.Printf(fmt.Sprintf("KafkaService.ProduceMessage() - Error: %v", err.Error()))
	}
	
	message := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder(receivedMessage)}
	producer.Input() <- message
	producer.AsyncClose() 
	
	log.Printf("KafkaService.ProduceMessage() - Successfully produced message!")
	
	return err;
}
