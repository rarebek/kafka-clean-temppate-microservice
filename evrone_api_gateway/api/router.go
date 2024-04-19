package api

import (
	"evrone_api_gateway/internal/usecase"
	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"evrone_api_gateway/api/handlers"
	v1 "evrone_api_gateway/api/handlers/v1"
	grpcClients "evrone_api_gateway/internal/infrastructure/grpc_service_client"
	"evrone_api_gateway/internal/pkg/config"
	//"evrone_api_gateway/internal/usecase/app_version"
	//"evrone_api_gateway/internal/usecase/event"
	//"evrone_api_gateway/internal/usecase/product"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Enforcer       *casbin.CachedEnforcer
	Service        grpcClients.ServiceClient
	//RefreshToken   product.RefreshToken
	BrokerProducer usecase.BrokerProducer
	//AppVersion     app_version.AppVersion
}

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) http.Handler {
	handleOption := &handlers.HandlerOption{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		//Enforcer:       option.Enforcer,
		Service: option.Service,
		//RefreshToken:   option.RefreshToken,
		//AppVersion:     option.AppVersion,
		BrokerProducer: option.BrokerProducer,
	}

	router := chi.NewRouter()
	router.Use(chimiddleware.RealIP, chimiddleware.Logger, chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(option.ContextTimeout))
	router.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/product", v1.NewProductHandler(handleOption))

	})

	// declare swagger api route
	router.Get("/swagger/*", httpSwagger.Handler())
	return router
}
