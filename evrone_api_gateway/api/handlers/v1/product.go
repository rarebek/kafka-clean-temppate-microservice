package v1

import (
	"encoding/json"
	_ "evrone_api_gateway/api/docs"
	"evrone_api_gateway/api/models"
	product_service "evrone_api_gateway/genproto/product"
	"evrone_api_gateway/internal/usecase"
	"fmt"
	"github.com/google/uuid"
	"net/http"

	errorapi "evrone_api_gateway/api/errors"
	"evrone_api_gateway/api/handlers"
	grpcClient "evrone_api_gateway/internal/infrastructure/grpc_service_client"
	"evrone_api_gateway/internal/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type productHandler struct {
	handlers.BaseHandler
	logger   *zap.Logger
	config   *config.Config
	service  grpcClient.ServiceClient
	producer usecase.BrokerProducer
	//enforcer *casbin.CachedEnforcer
}

func NewProductHandler(option *handlers.HandlerOption) http.Handler {
	handler := productHandler{
		logger:   option.Logger,
		config:   option.Config,
		service:  option.Service,
		producer: option.BrokerProducer,
		//enforcer: option.Enforcer,
	}

	//handler.Cache = option.Cache
	handler.Client = option.Service
	handler.Config = option.Config

	//policies := [][]string{
	//	{"investor", "/v1/content/categories", "GET"},
	//	{"investor", "/v1/content/chapters", "GET"},
	//	{"investor", "/v1/content/articles/{chapter_id}", "GET"},
	//	{"investor", "/v1/content/news", "GET"},
	//}
	//for _, policy := range policies {
	//	_, err := option.Enforcer.AddPolicy(policy)
	//	if err != nil {
	//		option.Logger.Error("error during investor enforcer add policies", zap.Error(err))
	//	}
	//}

	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		// auth
		//r.Use(middleware.Authorizer(option.Enforcer, option.Logger))

		// content
		r.Post("/add-product", handler.addProduct())
		r.Get("/get-product", handler.getProduct())
		r.Delete("/delete-product", handler.deleteProduct())
		r.Put("/update-product", handler.updateProduct())

	})
	return router
}

// @Summary     Product
// @Description Add a product to inventory
// @ID          add-product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       request body models.Product true "Enter product details"
// @Success     200 {object} models.Product
// @Failure     400 {object} errors.ErrResponse
// @Failure     500 {object} errors.ErrResponse
// @Router      /v1/product/add-product [post]
func (h *productHandler) addProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var (
			product models.Product
		)

		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			h.logger.Error("error decoding product", zap.Error(err))
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = h.producer.ProduceUserInfoToKafka(ctx, "1", &product_service.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			UnitPrice:   float32(product.UnitPrice),
		})
		if err != nil {
			h.logger.Error("kafka error", zap.Error(err))
			http.Error(w, "Kafka error", http.StatusInternalServerError)
			return
		}

		id := uuid.NewString()

		product.ID = id

		resp, err := h.service.ProductService().AddProduct(ctx, &product_service.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Name,
			CategoryID:  product.CategoryID,
			UnitPrice:   float32(product.UnitPrice),
		})

		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, resp)
	}
}

// @Summary     Product
// @Description Get one product details from inventory
// @ID          get-product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       id query string true "Product ID"
// @Success     200 {object} models.Product
// @Failure     400 {object} errors.ErrResponse
// @Failure     500 {object} errors.ErrResponse
// @Router      /v1/product/get-product [get]
func (h *productHandler) getProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.URL.Query().Get("id")
		if id == "" {
			render.Render(w, r, errorapi.Error(fmt.Errorf("id parameter is required")))
		}

		resp, err := h.service.ProductService().GetProduct(ctx, &product_service.IdRequest{ID: id})
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, resp)
	}
}

// @Summary     Product
// @Description Delete one product details from inventory
// @ID          delete-product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       id query string true "Product ID"
// @Success     200 {object} models.Product
// @Failure     400 {object} errors.ErrResponse
// @Failure     500 {object} errors.ErrResponse
// @Router      /v1/product/delete-product [delete]
func (h *productHandler) deleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.URL.Query().Get("id")
		if id == "" {
			render.Render(w, r, errorapi.Error(fmt.Errorf("id parameter is required")))
		}

		resp, err := h.service.ProductService().DeleteProduct(ctx, &product_service.IdRequest{ID: id})
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, resp)
	}
}

// @Summary     Product
// @Description Update a product to inventory
// @ID          update-product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       id query string true "Type ProductID"
// @Param       request body models.Product true "Enter product details"
// @Success     200 {object} models.Product
// @Failure     400 {object} errors.ErrResponse
// @Failure     500 {object} errors.ErrResponse
// @Router      /v1/product/update-product [put]
func (h *productHandler) updateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var (
			product models.Product
		)
		id := r.URL.Query().Get("id")

		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			h.logger.Error("error decoding product", zap.Error(err))
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		product.ID = id

		resp, err := h.service.ProductService().UpdateProduct(ctx, &product_service.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			UnitPrice:   float32(product.UnitPrice),
		})
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, resp)
	}
}
