package kafkaclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Producer struct {
	kafkaProducer *kafka.Producer
	paymentsTopic string
	logger        *zap.Logger
}

func New(brokers []string, paymentsTopic string, logger *zap.Logger) (*Producer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":   strings.Join(brokers, ","),
		"acks":                "all",
		"enable.idempotence":  true,
		"retries":             3,
		"delivery.timeout.ms": 10000,
		"request.timeout.ms":  5000,
	}

	kp, err := kafka.NewProducer(cfg)

	if err != nil {
		return nil, fmt.Errorf("kafka: new producer: %w", err)
	}

	p := &Producer{
		kafkaProducer: kp,
		paymentsTopic: paymentsTopic,
		logger:        logger,
	}

	go p.drainEvents()

	return p, nil
}

func (p *Producer) drainEvents() {
	for ev := range p.kafkaProducer.Events() {
		if kerr, ok := ev.(kafka.Error); ok {
			p.logger.Error("kafka producer event", zap.Error(kerr))
		}
	}
}

func (p *Producer) Publish(ctx context.Context, key string, payload []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.paymentsTopic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: payload,
	}

	deliveryChan := make(chan kafka.Event, 1)

	if err :=  p.kafkaProducer.Produce(message, deliveryChan); err != nil {
		return fmt.Errorf("kafka: produce: %w", err)
	}

	select {
	case ev := <-deliveryChan:
		m, ok := ev.(*kafka.Message)

		if !ok {
			return fmt.Errorf("kafka: unexpected delivery event type: %T", ev)
		}

		if m.TopicPartition.Error != nil {
			return fmt.Errorf("kafka: delivery: %w", m.TopicPartition.Error)
		}

		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) Close() error {
	remaining := p.kafkaProducer.Flush(5000)
	p.kafkaProducer.Close()

	if remaining > 0 {
		return fmt.Errorf("kafka: %d messages not flushed on close", remaining)
	}
	return nil
}
