package msgBroker

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/degarzonm/customer_leal_service/internal/config"
	"github.com/degarzonm/customer_leal_service/internal/domain"
	"github.com/degarzonm/customer_leal_service/internal/infrastructure/util"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// NewKafkaProducer creates and returns a new Kafka producer instance.
// It uses the global configuration to set up the producer with specific settings,
// such as returning successes and using a hash partitioner. If there is an error
// in creating the producer, it returns the error. Otherwise, it returns an instance
// of KafkaProducer that implements the domain.EventProducer interface.
func NewKafkaProducer() (domain.EventProducer, error) {
	cfg := config.GetConfig()
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{Producer: producer}, nil
}

// SendMessage sends a message to a specified Kafka topic.
//
// It generates a unique key for the message using the util.GenerateToken function.
// The message is marshaled into JSON format and sent as a sarama.ProducerMessage
// with the specified topic and the generated key. If the message is successfully
// sent, it logs the partition and offset of the message in the topic.
//
// Returns an error if the key generation, message marshaling, or message sending
// fails.

func (kp *KafkaProducer) SendMessage(topic string, message interface{}) error {
	key, err := util.GenerateToken()
	if err != nil {
		return err
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(messageBytes),
	}

	partition, offset, err := kp.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s [partition=%d, offset=%d]\n", topic, partition, offset)
	return nil
}
