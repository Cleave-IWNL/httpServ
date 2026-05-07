package kafkaclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"httpServ/internal/model"

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

func (p *Producer) PublishPaymentCreated(ctx context.Context, event model.PaymentCreatedEvent) error {
	payload, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("kafka: marshal event: %w", err)
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.paymentsTopic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(event.PaymentID),
		Value: payload,
	}

	deliveryChan := make(chan kafka.Event, 1)

	if err := p.kafkaProducer.Produce(message, deliveryChan); err != nil {
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

func (p *Producer) Close() {
	remeaning := p.kafkaProducer.Flush(5000)
	if remeaning > 0 {
		p.logger.Warn("kafkaa priduces: messages not flushed on close",
			zap.Int("remeaning", remeaning))
	}

	p.kafkaProducer.Close()
}
