package kafka

import (
	"context"
	"errors"
	"evrone_api_gateway/internal/usecase"
	"fmt"
	"io"

	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	minBytes = 1e3 // ~1KB
	maxBytes = 1e6 // ~1MB
)

type consumer struct {
	logger  *zap.Logger
	readers []*kafka.Reader
	cfg     []usecase.ConsumerConfig
}

func NewConsumer(logger *zap.Logger) usecase.BrokerConsumer {
	return &consumer{
		logger:  logger,
		cfg:     make([]usecase.ConsumerConfig, 0),
		readers: make([]*kafka.Reader, 0),
	}
}

func (c *consumer) Run() error {
	for _, cfg := range c.cfg {
		var reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:     cfg.GetBrokers(),
			GroupID:     cfg.GetGroupID(),
			Topic:       cfg.GetTopic(),
			MinBytes:    minBytes,
			MaxBytes:    maxBytes,
			Logger:      nil,
			ErrorLogger: nil,
		})

		go func(reader *kafka.Reader, cfg usecase.ConsumerConfig) {
			runReader(reader, c.logger, cfg)
		}(reader, cfg)

		c.readers = append(c.readers, reader)
	}

	return nil
}

func runReader(r *kafka.Reader, logger *zap.Logger, cfg usecase.ConsumerConfig) {
	var (
		topic           = cfg.GetTopic()
		failMessageRead = func(topic string) string {
			return fmt.Sprintf("consumer failed to read a message in %s", topic)
		}
	)

	for {
		ctx := context.Background()
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			logger.Error(failMessageRead(topic), zap.Error(err))
			continue
		}

		if err = r.CommitMessages(ctx, msg); err != nil {
			logger.Error("consumer failed to commit message", zap.Error(err), zap.String("topic", r.Config().Topic))
			continue
		}

		if err = cfg.GetHandler().Handle(ctx, msg.Key, msg.Value); err != nil {
			logger.Error("handler failed to process message", zap.Error(err), zap.String("topic", r.Config().Topic), zap.ByteString("value", msg.Value))
		}
	}
}

func (c *consumer) RegisterConsumer(cfg usecase.ConsumerConfig) {
	c.cfg = append(c.cfg, cfg)
}

func (c *consumer) Close() {
	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			c.logger.Error("consumer reader close", zap.Error(err))
		}
	}
}

type consumerConfig struct {
	Brokers []string
	Topic   string
	GroupId string
	Handler usecase.ConsumerHandler
}

func NewCustomerConsumerConfig(brokers []string, topic, groupId string, handler usecase.ConsumerHandler) usecase.ConsumerConfig {
	return &consumerConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupId: groupId,
		Handler: handler,
	}
}

func (c *consumerConfig) GetBrokers() []string {
	return c.Brokers
}

func (c *consumerConfig) GetTopic() string {
	return c.Topic
}

func (c *consumerConfig) GetGroupID() string {
	return c.GroupId
}

func (c *consumerConfig) GetHandler() usecase.ConsumerHandler {
	return c.Handler
}
