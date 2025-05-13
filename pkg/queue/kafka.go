package queue

import (
	"crypto/tls"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
)

func ConnectProducer(brokerUrls []string, apiKey, secret string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	if apiKey != "" && secret != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = apiKey
		config.Net.SASL.Password = secret
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
	}

	producer, err := sarama.NewSyncProducer(brokerUrls, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func PushMessageWithKeyToQueue(brokerUrls []string, apiKey, secret, topic, key string, message []byte) error {
	producer, err := ConnectProducer(brokerUrls, apiKey, secret)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("message is stored in topic(%s)/partision(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}

func ConnectConsumer(brokerUrls []string, apiKey, secret string) (sarama.Consumer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 3

	if apiKey != "" && secret != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = apiKey
		config.Net.SASL.Password = secret
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
	}

	consumer, err := sarama.NewConsumer(brokerUrls, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func DecodeMessage(obj any, value []byte) error {
	if err := json.Unmarshal(value, &obj); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		return err
	}

	return nil
}
