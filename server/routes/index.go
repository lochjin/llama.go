package routes

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

//go:embed assets/index.html.gz
var indexHTMLGzip []byte //For more source code details, please go to https://github.com/Qitmeer/llamago-webui

func (s *API) IndexHandler(c *gin.Context) {
	c.Header("Vary", "Accept-Encoding")
	c.Header("Cache-Control", "public, max-age=600")
	c.Header("Content-Type", "text/html; charset=utf-8")

	acceptEncoding := c.GetHeader("Accept-Encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		c.Header("Content-Encoding", "gzip")
		if c.Request.Method == "HEAD" {
			c.Status(200)
			return
		}
		c.Data(200, "text/html; charset=utf-8", indexHTMLGzip)
		return
	}

	zr, err := gzip.NewReader(bytes.NewReader(indexHTMLGzip))
	if err != nil {
		c.String(500, "failed to decode gzip")
		return
	}
	defer zr.Close()

	if c.Request.Method == "HEAD" {
		c.Status(200)
		return
	}

	data, err := io.ReadAll(zr)
	if err != nil {
		c.String(500, "failed to read content")
		return
	}
	c.Data(200, "text/html; charset=utf-8", data)
}
