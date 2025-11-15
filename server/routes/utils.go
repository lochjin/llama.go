package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Qitmeer/llama.go/api"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
)

type ImageData struct {
	Data []byte `json:"data"`
	ID   int    `json:"id"`
}

func waitForStream(c *gin.Context, ch chan any) {
	c.Header("Content-Type", "application/json")
	var latest api.ProgressResponse
	for resp := range ch {
		switch r := resp.(type) {
		case api.ProgressResponse:
			latest = r
		case gin.H:
			status, ok := r["status"].(int)
			if !ok {
				status = http.StatusInternalServerError
			}
			errorMsg, ok := r["error"].(string)
			if !ok {
				errorMsg = "unknown error"
			}
			c.JSON(status, gin.H{"error": errorMsg})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown message type"})
			return
		}
	}

	c.JSON(http.StatusOK, latest)
}

func streamHandler(c *gin.Context, ch chan any) {
	accept := c.GetHeader("Accept")
	if accept == "application/x-ndjson" {
		// NDJSON
		c.Header("Content-Type", "application/x-ndjson")

		c.Stream(func(w io.Writer) bool {
			val, ok := <-ch
			if !ok {
				return false
			}

			bts, ok := val.(string)
			if !ok {
				jbts, err := json.Marshal(val)
				if err != nil {
					log.Warn("NDJSON marshal error", "error", err)
					return false
				} else {
					bts = string(jbts)
				}
			}
			bts += "\n"
			if _, err := w.Write([]byte(bts)); err != nil {
				log.Warn("NDJSON write error:", err)
				return false
			}

			return true
		})
	} else if accept == "text/event-stream" {
		// SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")

		c.Stream(func(w io.Writer) bool {
			val, ok := <-ch
			if !ok {
				return false
			}
			bts, ok := val.(string)
			if !ok {
				log.Warn("SSE marshal error", "error", val)
				return false
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", bts); err != nil {
				log.Warn("SSE write error:", err)
				return false
			}
			return true
		})
	} else {
		c.Stream(func(w io.Writer) bool {
			val, ok := <-ch
			if !ok {
				return false
			}
			bts, ok := val.(string)
			if !ok {
				log.Warn("default marshal error", "error", val)
				return false
			}
			if _, err := w.Write([]byte(bts)); err != nil {
				log.Warn("default write error:", err)
				return false
			}
			return true
		})
	}
}
