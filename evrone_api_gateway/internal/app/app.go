package app

import (
	"context"
	"evrone_api_gateway/internal/infrastructure/kafka"
	"evrone_api_gateway/internal/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"evrone_api_gateway/api"
	grpcService "evrone_api_gateway/internal/infrastructure/grpc_service_client"
	"evrone_api_gateway/internal/pkg/config"
	"evrone_api_gateway/internal/pkg/logger"

	"evrone_api_gateway/internal/pkg/postgres"
	"evrone_api_gateway/internal/pkg/redis"
	//"evrone_api_gateway/internal/usecase/app_version"
	//"evrone_api_gateway/internal/usecase/event"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	RedisDB        *redis.RedisDB
	server         *http.Server
	Enforcer       *casbin.CachedEnforcer
	Clients        grpcService.ServiceClient
	ShutdownOTLP   func() error
	BrokerProducer usecase.BrokerProducer
	//appVersion     app_version.AppVersion
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}
	//
	//// kafka producer init
	kafkaProducer := kafka.NewProducer(&cfg, logger)

	// postgres init
	db, err := postgres.New(&cfg)
	if err != nil {
		return nil, err
	}

	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}

	//// otlp collector init
	//shutdownOTLP, err := otlp.InitOTLPProvider(&cfg)
	//if err != nil {
	//	return nil, err
	//}

	// initialization enforcer
	//enforcer, err := policy.NewCachedEnforcer(&cfg, logger)
	//if err != nil {
	//	return nil, err
	//}

	//enforcer.SetCache(policy.NewCache(&redisdb.Client))
	//
	//var (
	//	contextTimeout time.Duration
	//)
	//
	//// context timeout initialization
	//contextTimeout, err = time.ParseDuration(cfg.Context.Timeout)
	//if err != nil {
	//	return nil, err
	//}
	////
	//appVersionRepo := postgresql.NewAppVersionRepo(db)
	//
	//appVersionUseCase := app_version.NewAppVersionService(contextTimeout, appVersionRepo)

	return &App{
		Config:  &cfg,
		Logger:  logger,
		DB:      db,
		RedisDB: redisdb,
		//Enforcer: enforcer,
		BrokerProducer: kafkaProducer,
		//ShutdownOTLP:   shutdownOTLP,
		//appVersion:     appVersionUseCase,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	clients, err := grpcService.New(a.Config)
	if err != nil {
		return err
	}
	a.Clients = clients
	//
	//// initialize cache
	//cache := redisrepo.NewCache(a.RedisDB)

	//tokenRepo := postgresql.NewRefreshTokenRepo(a.DB)
	//
	//initialize token service
	//refreshTokenService := product.NewRefreshTokenService(contextTimeout, tokenRepo)

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		//Cache:          cache,
		Enforcer: a.Enforcer,
		//RefreshToken:   refreshTokenService,
		Service:        clients,
		BrokerProducer: a.BrokerProducer,
		//AppVersion:     a.appVersion,
	})
	//if err = a.Enforcer.LoadPolicy(); err != nil {
	//	return fmt.Errorf("error during enforcer load policy: %w", err)
	//}

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// close grpc connections
	a.Clients.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
