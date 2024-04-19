package usecase

import (
	"context"
	product_service "evrone_api_gateway/genproto/product"
)

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupID() string
	GetHandler() ConsumerHandler
}

type ConsumerHandler interface {
	Handle(ctx context.Context, key, value []byte) error
}

type BrokerConsumer interface {
	Run() error
	RegisterConsumer(cfg ConsumerConfig)
	Close()
}

type BrokerProducer interface {
	ProduceUserInfoToKafka(ctx context.Context, key string, body *product_service.Product) error
	Close()
}
