package kafka

import (
	"context"
	"encoding/json"
	product_service "evrone_api_gateway/genproto/product"
	configpkg "evrone_api_gateway/internal/pkg/config"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type producer struct {
	logger         *zap.Logger
	productService *kafka.Writer
}

func NewProducer(config *configpkg.Config, logger *zap.Logger) *producer {
	return &producer{
		logger: logger,
		productService: &kafka.Writer{
			Addr:                   kafka.TCP(config.Kafka.Address...),
			Topic:                  config.Kafka.Topic.ProductService,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireAll,
			AllowAutoTopicCreation: true,
			Async:                  true,
			Completion: func(messages []kafka.Message, err error) {
				if err != nil {
					logger.Error("kafka product created", zap.Error(err))
				}
				for _, message := range messages {
					logger.Sugar().Info(
						"kafka investmentCreated message",
						zap.Int("partition", message.Partition),
						zap.Int64("offset", message.Offset),
						zap.String("key", string(message.Key)),
						zap.String("value", string(message.Value)),
					)
				}
			},
		},
	}
}

func (p *producer) buildMessage(key string, value []byte) kafka.Message {
	return kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
}

func (p *producer) ProduceUserInfoToKafka(ctx context.Context, key string, body *product_service.Product) error {
	byteData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	fmt.Println(body)
	message := p.buildMessage(key, byteData)

	return p.productService.WriteMessages(ctx, message)
}

func (p *producer) Close() {
	if err := p.productService.Close(); err != nil {
		p.logger.Error("error during close writer investmentCreated", zap.Error(err))
	}
}
