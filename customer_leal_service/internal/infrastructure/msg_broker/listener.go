package msgBroker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/degarzonm/customer_leal_service/internal/application"
	"github.com/degarzonm/customer_leal_service/internal/config"
	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type KafkaListener struct {
	consumerGroup sarama.ConsumerGroup
	appService    *application.AppService
}

// NewKafkaListener creates a new Kafka listener instance that consumes
// messages from the topic specified in the configuration, and passes them
// to the provided application service for processing.
// The returned listener instance is ready to be used with the Listen
// method to start consuming messages.
func NewKafkaListener(appService *application.AppService) (*KafkaListener, error) {
	cfg := config.GetConfig()
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.CustomerGroup, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &KafkaListener{
		consumerGroup: consumerGroup,
		appService:    appService,
	}, nil
}

// Listen starts consuming messages from the MsgApplyPointsTopic specified in the configuration.
//
// The method will loop indefinitely, logging any errors that occur while
// consuming messages. It returns an error if the consumer group fails
// to consume from the topics.

func (kl *KafkaListener) Listen() error {
	cfg := config.GetConfig()
	topics := []string{cfg.MsgApplyPointsTopic}

	for {
		if err := kl.consumerGroup.Consume(context.Background(), topics, kl); err != nil {
			log.Printf("Error consuming Kafka topics: %v", err)
			return err
		}
	}
}

func (kl *KafkaListener) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (kl *KafkaListener) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from the Kafka topic.
//
// The method will loop indefinitely, logging any errors that occur while
// consuming messages. The method will unmarshal purchase messages from the
// MsgApplyPointsTopic topic, and pass them to the application layer to be
// processed. Any errors that occur while processing the purchase will be
// logged.
//
// If the method encounters an unhandled topic, it will log a message indicating
// this.
func (kl *KafkaListener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	cfg := config.GetConfig()

	for message := range claim.Messages() {
		log.Println("Processing message from topic:", message.Topic)
		switch message.Topic {
		case cfg.MsgApplyPointsTopic:
			log.Println("topic apply points with message", message.Value)
			var points domain.LealPointsApply
			if err := json.Unmarshal(message.Value, &points); err != nil {
				log.Printf("Error unmarshalling purchase: %v", err)
				continue
			}
			log.Println("unmarshalled points", points)
			if err := kl.appService.ProcessApplyPointsEvent(points); err != nil {
				log.Printf("Error processing purchase: %v", err)
			}
		default:
			log.Printf("Unhandled topic: %s", message.Topic)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
