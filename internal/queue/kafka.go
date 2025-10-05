package queue

import (
	"context"
	"encoding/json"
	"time"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Publisher interface {
	Publish(ctx context.Context, msg any) error
	Close() error
}

type KafkaPublisher struct{ writer *kafka.Writer }

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{writer: &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
		BatchTimeout: 50 * time.Millisecond,
	}}
}

func (p *KafkaPublisher) Publish(ctx context.Context, msg any) error {
	b, err := json.Marshal(msg)
	if err != nil { return err }
	if err := p.writer.WriteMessages(ctx, kafka.Message{Value: b}); err != nil {
		log.Error().Err(err).Msg("kafka write failed")
		return err
	}
	return nil
}

func (p *KafkaPublisher) Close() error {
	if p.writer != nil { return p.writer.Close() }
	return nil
}
