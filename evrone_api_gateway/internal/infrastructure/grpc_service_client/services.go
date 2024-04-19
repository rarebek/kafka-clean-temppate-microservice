package grpc_service_clients

import (
	product_service "evrone_api_gateway/genproto/product"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"evrone_api_gateway/internal/pkg/config"
)

type ServiceClient interface {
	ProductService() product_service.ProductServiceClient
	Close()
}

type serviceClient struct {
	connections    []*grpc.ClientConn
	productService product_service.ProductServiceClient
}

func New(cfg *config.Config) (ServiceClient, error) {
	connProductService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ContentService.Host, cfg.ContentService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClient{
		productService: product_service.NewProductServiceClient(connProductService),
		connections: []*grpc.ClientConn{
			connProductService,
		},
	}, nil
}

func (s *serviceClient) ProductService() product_service.ProductServiceClient {
	return s.productService
}

func (s *serviceClient) Close() {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			// should be replaced by logger soon
			fmt.Printf("error while closing grpc connection: %v", err)
		}
	}
}
