package handlers

import (
	"evrone_api_gateway/internal/usecase"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	"evrone_api_gateway/api/middleware"
	grpcClients "evrone_api_gateway/internal/infrastructure/grpc_service_client"
	"evrone_api_gateway/internal/pkg/config"
	//"evrone_api_gateway/internal/pkg/otlp"
	//appV "evrone_api_gateway/internal/usecase/app_version"
	//"evrone_api_gateway/internal/usecase/event"
	//"evrone_api_gateway/internal/usecase/product"
)

const (
	InvestorToken = "investor"
)

type HandlerOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Enforcer       *casbin.CachedEnforcer
	//Cache          redis.Cache
	Service grpcClients.ServiceClient
	//RefreshToken   product.RefreshToken
	//AppVersion     appV.AppVersion
	BrokerProducer usecase.BrokerProducer
}

type BaseHandler struct {
	//Cache  redis.Cache
	Config *config.Config
	Client grpcClients.ServiceClient
}

func (h *BaseHandler) GetAuthData(ctx context.Context) (map[string]string, bool) {
	//// tracing
	//ctx, span := otlp.Start(ctx, "handler", "GetAuthData")
	//defer span.End()

	data, ok := ctx.Value(middleware.RequestAuthCtx).(map[string]string)
	return data, ok
}
