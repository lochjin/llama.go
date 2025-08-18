package server

import (
	"errors"
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/version"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ollama/ollama/openai"
	"github.com/ollama/ollama/template"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
	"sync"
)

type Service struct {
	ctx  *cli.Context
	cfg  *config.Config
	tmpl *template.Template

	addr net.Addr
	srvr *http.Server

	wg sync.WaitGroup
}

func New(ctx *cli.Context, cfg *config.Config) *Service {
	log.Info("New Server ...")
	ser := Service{ctx: ctx, cfg: cfg}
	return &ser
}

func (s *Service) Start() error {
	log.Info("Start Server...")
	tmpl, err := template.Parse("{{- range .Messages }}<|im_start|>{{ .Role }}\n{{ .Content }}<|im_end|>\n{{ end }}<|im_start|>assistant")
	if err != nil {
		return err
	}
	s.tmpl = tmpl

	ln, err := net.Listen("tcp", s.cfg.HostURL().Host)
	if err != nil {
		return err
	}
	s.addr = ln.Addr()

	err = s.GenerateRoutes()
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Listening on %s (version %s)", ln.Addr(), version.String()))
	s.srvr = &http.Server{
		Handler: nil,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err = s.srvr.Serve(ln)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error(err.Error())
		}
	}()

	return nil
}

func (s *Service) GenerateRoutes() error {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowWildcard = true
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowHeaders = []string{
		"Authorization",
		"Content-Type",
		"User-Agent",
		"Accept",
		"X-Requested-With",

		// OpenAI compatibility headers
		"OpenAI-Beta",
		"x-stainless-arch",
		"x-stainless-async",
		"x-stainless-custom-poll-interval",
		"x-stainless-helper-method",
		"x-stainless-lang",
		"x-stainless-os",
		"x-stainless-package-version",
		"x-stainless-poll-helper",
		"x-stainless-retry-count",
		"x-stainless-runtime",
		"x-stainless-runtime-version",
		"x-stainless-timeout",
	}
	corsConfig.AllowOrigins = s.cfg.AllowedOrigins()

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.HandleMethodNotAllowed = true
	r.Use(
		cors.New(corsConfig),
		allowedHostsMiddleware(s.addr),
	)

	// General
	r.HEAD("/", func(c *gin.Context) { c.String(http.StatusOK, "Llamago is running") })
	r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Llamago is running") })
	r.HEAD("/api/version", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"version": version.String()}) })
	r.GET("/api/version", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"version": version.String()}) })

	r.HEAD("/api/tags", s.ListHandler)
	r.GET("/api/tags", s.ListHandler)
	r.HEAD("/api/models", s.ListHandler)
	r.GET("/api/models", s.ListHandler)

	r.POST("/api/show", s.ShowHandler)

	r.GET("/api/ps", s.PsHandler)
	r.POST("/api/generate", s.GenerateHandler)
	r.POST("/api/chat", s.ChatHandler)
	r.POST("/api/embed", s.EmbedHandler)
	r.POST("/api/embeddings", s.EmbeddingsHandler)

	// Inference (OpenAI compatibility)
	r.POST("/v1/chat/completions", openai.ChatMiddleware(), s.ChatHandler)
	r.POST("/v1/completions", openai.CompletionsMiddleware(), s.GenerateHandler)
	r.POST("/v1/embeddings", openai.EmbeddingsMiddleware(), s.EmbedHandler)
	r.GET("/v1/models", openai.ListMiddleware(), s.ListHandler)
	r.GET("/v1/models/:model", openai.RetrieveMiddleware(), s.ShowHandler)

	http.Handle("/", r)
	return nil
}

func (s *Service) Stop() error {
	log.Info("Stop Server...")

	var err error
	if s.srvr != nil {
		err = s.srvr.Close()
	}
	s.wg.Wait()
	return err
}
