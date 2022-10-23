package connectors

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type Conf struct {
	KafkaBrokers       string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaTopic         string `envconfig:"KAFKA_TOPICS" required:"true"`
	KafkaConsumerGroup string `envconfig:"KAFKA_CONSUMER_GROUP" required:"true"`

	KafkaTLSCACertFile string `envconfig:"KAFKA_TLS_CA_CERT_FILE" default:"/secrets/ca_cert.pem"`
	KafkaTLSClientCert string `envconfig:"KAFKA_TLS_CLIENT_CERT" default:"/secrets/client_cert.pem"`
	KafkaTLSClientKey  string `envconfig:"KAFKA_TLS_CLIENT_KEY" default:"/secrets/client_key.pem"`
	KafkaTLSEnabled    bool   `envconfig:"KAFKA_TLS_ENABLED" default:"false"`
}

func Load() (Conf, error) {
	c := Conf{}
	if err := envconfig.Process("", &c); err != nil {
		logrus.Errorf(fmt.Sprintf("load env error: %v", err))
	}
	return c, nil
}

type KafkaProducer struct {
	Config       Conf
	Producer     sarama.SyncProducer
	Reconnecting bool
	IsClosed     bool
}

func (k *KafkaProducer) Init() error {
	config, err := Load()
	if err != nil {
		return err
	}
	k.Config = config
	k.Producer, err = newKafkaProducer(config)
	return err
}

func (k *KafkaProducer) Reconnect() {
	if k.Reconnecting {
		return
	}
	defer func() {
		k.Reconnecting = false
	}()
	err := k.Init()
	if err != nil {
		logrus.Errorf("Error in reconnect kafka producer! Sleep 10s!")
	}
	time.Sleep(10 * time.Second)
	k.Reconnect()
}

func (k *KafkaProducer) Close() error {
	k.IsClosed = true
	if k.Producer == nil {
		return errors.New("Producer is nil!")
	}
	if err := k.Producer.Close(); err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

func (k *KafkaProducer) SendMessageData(topic, event string) {
	if k.IsClosed {
		logrus.Errorf("Produce closed!")
		return
	}
	if k.Producer == nil {
		logrus.Errorf("Invalid producer! Retry Reconnect")
		k.Reconnect()
		return
	}
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(event)}
	partition, offset, err := k.Producer.SendMessage(msg)
	if err != nil {
		logrus.Errorf(fmt.Sprintf("FAILED to send message: %s\n", err))
	} else {
		logrus.Infof(fmt.Sprintf("> message: %v sent to partition %d at offset %d\n", event, partition, offset))
	}
}

func (k *KafkaProducer) Notification(topic string, payload interface{}) {
	producer := KafkaProducer{}
	if err := producer.Init(); err != nil {
		logrus.Errorf("Error in Init Kafka producer, err: %v", err.Error())
		return
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logrus.Errorf("Error in Init Kafka producer, err: %v", err.Error())
			return
		}
	}()
	dataBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	producer.SendMessageData(topic, string(dataBytes))
}

func newKafkaProducer(config Conf) (sarama.SyncProducer, error) {
	configKafka := sarama.NewConfig()
	configKafka.Version = sarama.V1_0_0_0

	configKafka.Producer.Return.Successes = true
	if config.KafkaTLSEnabled {
		tlsConfig, err := newTLSConfig(config.KafkaTLSClientCert, config.KafkaTLSClientKey, config.KafkaTLSCACertFile)
		if err != nil {
			logrus.Errorf(fmt.Sprintf("Setup kafka TLS error: %v", err))
			return nil, err
		}
		tlsConfig.InsecureSkipVerify = true
		configKafka.Net.TLS.Enable = true
		configKafka.Net.TLS.Config = tlsConfig
	}

	producer, err := sarama.NewSyncProducer(strings.Split(config.KafkaBrokers, ","), configKafka)
	if err != nil {
		logrus.Errorf(fmt.Sprintf("Setup kafka error: %v", err))
		return nil, err
	}

	return producer, nil
}

func newTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool
	return &tlsConfig, err
}
