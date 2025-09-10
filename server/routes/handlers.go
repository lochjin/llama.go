package routes

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	qapi "github.com/Qitmeer/llama.go/api"
	"github.com/Qitmeer/llama.go/model"
	"github.com/Qitmeer/llama.go/version"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"github.com/ollama/ollama/api"
	omodel "github.com/ollama/ollama/types/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func (s *API) PsHandler(c *gin.Context) {
	models := []api.ProcessModelResponse{}
	slices.SortStableFunc(models, func(i, j api.ProcessModelResponse) int {
		// longest duration remaining listed first
		return cmp.Compare(j.ExpiresAt.Unix(), i.ExpiresAt.Unix())
	})

	c.JSON(http.StatusOK, api.ProcessResponse{Models: models})
}

func (s *API) GenerateHandler(c *gin.Context) {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bodyStr := string(bodyBytes)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req api.GenerateRequest
	if err := c.ShouldBindJSON(&req); errors.Is(err, io.EOF) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Model) > 0 {
		if s.cfg.Model != req.Model {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("model '%s' not found", req.Model)})
			return
		}
	}

	// expire the runner
	if req.Prompt == "" && req.KeepAlive != nil && int(req.KeepAlive.Seconds()) == 0 {
		c.JSON(http.StatusOK, api.GenerateResponse{
			Model:      req.Model,
			CreatedAt:  time.Now().UTC(),
			Response:   "Not currently supported",
			Done:       true,
			DoneReason: "unload",
		})
		return
	}

	if req.Raw && (req.Template != "" || req.System != "" || len(req.Context) > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "raw mode does not support template, system, or context"})
		return
	}

	caps := []omodel.Capability{omodel.CapabilityCompletion}
	if req.Suffix != "" {
		caps = append(caps, omodel.CapabilityInsert)
	}
	if req.Think != nil && *req.Think {
		caps = append(caps, omodel.CapabilityThinking)
		// TODO(drifkin): consider adding a warning if it's false and the model
		// doesn't support thinking. It's not strictly required, but it can be a
		// hint that the user is on an older qwen3/r1 model that doesn't have an
		// updated template supporting thinking
	}

	// load the model
	if req.Prompt == "" {
		c.JSON(http.StatusOK, api.GenerateResponse{
			Model:      req.Model,
			CreatedAt:  time.Now().UTC(),
			Done:       true,
			DoneReason: "load",
		})
		return
	}

	content, err := wrapper.LlamaGenerate(bodyStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var ret map[string]interface{}
	if err := json.Unmarshal([]byte(content), &ret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid json"})
		return
	}
	c.JSON(http.StatusOK, ret)
}

func (s *API) ChatHandler(c *gin.Context) {
	checkpointStart := time.Now()
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bodyStr := string(bodyBytes)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req api.ChatRequest
	if err := c.ShouldBindJSON(&req); errors.Is(err, io.EOF) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// expire the runner
	if len(req.Messages) == 0 && req.KeepAlive != nil && int(req.KeepAlive.Seconds()) == 0 {
		c.JSON(http.StatusOK, api.ChatResponse{
			Model:      req.Model,
			CreatedAt:  time.Now().UTC(),
			Message:    api.Message{Role: "assistant", Content: "Not currently supported"},
			Done:       true,
			DoneReason: "unload",
		})
		return
	}

	caps := []omodel.Capability{omodel.CapabilityCompletion}
	if len(req.Tools) > 0 {
		caps = append(caps, omodel.CapabilityTools)
	}
	if req.Think != nil && *req.Think {
		caps = append(caps, omodel.CapabilityThinking)
	}

	checkpointLoaded := time.Now()

	if len(req.Messages) == 0 {
		c.JSON(http.StatusOK, api.ChatResponse{
			Model:      req.Model,
			CreatedAt:  time.Now().UTC(),
			Message:    api.Message{Role: "assistant"},
			Done:       true,
			DoneReason: "load",
		})
		return
	}

	content, err := wrapper.LlamaChat(bodyStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := api.ChatResponse{
		Model:     req.Model,
		CreatedAt: time.Now().UTC(),
		Message:   api.Message{Role: "assistant", Content: content},
		Done:      true,
	}

	res.TotalDuration = time.Since(checkpointStart)
	res.LoadDuration = checkpointLoaded.Sub(checkpointStart)
	if req.Stream == nil && !*req.Stream {
		c.JSON(http.StatusOK, res)
		return
	}

	c.Header("Content-Type", "application/x-ndjson")
	c.Stream(func(w io.Writer) bool {
		bts, err := json.Marshal(res)
		if err != nil {
			log.Info(fmt.Sprintf("streamResponse: json.Marshal failed with %s", err))
			return false
		}

		// Delineate chunks with new-line delimiter
		bts = append(bts, '\n')
		if _, err := w.Write(bts); err != nil {
			log.Info(fmt.Sprintf("streamResponse: w.Write failed with %s", err))
			return false
		}
		return true
	})
}

func (s *API) EmbedHandler(c *gin.Context) {
	checkpointStart := time.Now()
	var req api.EmbedRequest
	err := c.ShouldBindJSON(&req)
	switch {
	case errors.Is(err, io.EOF):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	case err != nil:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input []string

	switch i := req.Input.(type) {
	case string:
		if len(i) > 0 {
			input = append(input, i)
		}
	case []any:
		for _, v := range i {
			if _, ok := v.(string); !ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input type"})
				return
			}
			input = append(input, v.(string))
		}
	default:
		if req.Input != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input type"})
			return
		}
	}

	checkpointLoaded := time.Now()

	if len(input) == 0 {
		c.JSON(http.StatusOK, api.EmbedResponse{Model: req.Model, Embeddings: [][]float32{}})
		return
	}

	prompts := ""
	for k, i := range input {
		if k > 0 {
			prompts += s.cfg.EmbdSeparator
		}
		prompts += i
	}

	ret, err := wrapper.LlamaEmbedding(s.cfg, s.cfg.ModelPath(), prompts, "array")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": strings.TrimSpace(err.Error())})
		return
	}
	var embeddings [][]float32
	err = json.Unmarshal([]byte(ret), &embeddings)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": strings.TrimSpace(err.Error())})
		return
	}
	if len(embeddings) != len(input) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%d != %d", len(embeddings), len(input))})
		return
	}
	resp := api.EmbedResponse{
		Model:           req.Model,
		Embeddings:      embeddings,
		TotalDuration:   time.Since(checkpointStart),
		LoadDuration:    checkpointLoaded.Sub(checkpointStart),
		PromptEvalCount: len(input),
	}
	c.JSON(http.StatusOK, resp)
}

func (s *API) EmbeddingsHandler(c *gin.Context) {
	var req api.EmbeddingRequest
	if err := c.ShouldBindJSON(&req); errors.Is(err, io.EOF) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// an empty request loads the model
	if req.Prompt == "" {
		c.JSON(http.StatusOK, api.EmbeddingResponse{Embedding: []float64{}})
		return
	}

	ret, err := wrapper.LlamaEmbedding(s.cfg, s.cfg.ModelPath(), req.Prompt, "array")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": strings.TrimSpace(err.Error())})
		return
	}
	var embeddings [][]float64
	err = json.Unmarshal([]byte(ret), &embeddings)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": strings.TrimSpace(err.Error())})
		return
	}
	resp := api.EmbeddingResponse{
		Embedding: embeddings[0],
	}
	c.JSON(http.StatusOK, resp)
}

func (s *API) ListHandler(c *gin.Context) {
	models := []api.ListModelResponse{}

	entries, err := os.ReadDir(s.cfg.ModelDir)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != model.EXT {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		models = append(models, api.ListModelResponse{
			Model:      entry.Name(),
			Name:       entry.Name(),
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
			Details: api.ModelDetails{
				Format: model.EXT[1:],
			},
		})
	}

	slices.SortStableFunc(models, func(i, j api.ListModelResponse) int {
		// most recently modified first
		return cmp.Compare(j.ModifiedAt.Unix(), i.ModifiedAt.Unix())
	})
	c.JSON(http.StatusOK, api.ListResponse{Models: models})
}

func (s *API) ShowHandler(c *gin.Context) {
	var req api.ShowRequest
	err := c.ShouldBindJSON(&req)
	switch {
	case errors.Is(err, io.EOF):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	case err != nil:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Model != "" {
		// noop
	} else if req.Name != "" {
		req.Model = req.Name
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "model is required"})
		return
	}

	entries, err := os.ReadDir(s.cfg.ModelDir)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != model.EXT {
			continue
		}
		if entry.Name() != req.Model {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		resp := &api.ShowResponse{
			Modelfile: info.Name(),
			Details: api.ModelDetails{
				Format: model.EXT[1:],
			},
			ModifiedAt: info.ModTime(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("model '%s' not found", req.Model)})
}

func (s *API) PropsHandler(c *gin.Context) {
	prep := qapi.PropsResponse{
		ModelPath: s.cfg.ModelPath(),
		BuildInfo: version.String(),
		NCtx:      int64(s.cfg.CtxSize),
		Modalities: qapi.Modalities{
			Vision: false,
			Audio:  false,
		},
	}
	c.JSON(http.StatusOK, prep)
}
