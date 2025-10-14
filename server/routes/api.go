package routes

import (
	"github.com/Qitmeer/llama.go/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
)

type API struct {
	cfg *config.Config
}

func New(cfg *config.Config) *API {
	log.Info("New API ...")
	ser := API{cfg: cfg}
	return &ser
}

func (s *API) Start() error {
	return nil
}

func (s *API) Setup(r *gin.Engine) {
	// General
	r.HEAD("/health", s.HealthHandler)
	r.GET("/health", s.HealthHandler)
	r.HEAD("/api/version", s.VersionHandler)
	r.GET("/api/version", s.VersionHandler)

	r.POST("/api/pull", s.PullHandler)
	r.HEAD("/api/tags", s.ListHandler)
	r.GET("/api/tags", s.ListHandler)
	r.HEAD("/api/models", s.ListHandler)
	r.GET("/api/models", s.ListHandler)
	r.POST("/api/show", s.ShowHandler)
	r.GET("/api/ps", s.PsHandler)
	r.GET("/props", s.PropsHandler)
	r.POST("/props", s.PropsChangeHandler)
	r.GET("/slots", s.SlotsHandler)

	r.POST("/api/generate", s.GenerateHandler)
	r.POST("/api/chat", s.ChatHandler)
	r.POST("/api/embed", s.EmbedHandler)
	r.POST("/api/embeddings", s.EmbeddingsHandler)

	// Inference (OpenAI compatibility)
	r.POST("/v1/completions", s.GenerateHandler)
	r.POST("/v1/chat/completions", s.ChatHandler)

	r.POST("/v1/embeddings", EmbeddingsMiddleware(), s.EmbedHandler)
	r.GET("/v1/models", ListMiddleware(), s.ListHandler)
	r.GET("/v1/models/:model", RetrieveMiddleware(), s.ShowHandler)

	// webui index
	r.GET("/", s.IndexHandler)
	r.HEAD("/", s.IndexHandler)
	r.GET("/index.html", s.IndexHandler)
	r.HEAD("/index.html", s.IndexHandler)
}
