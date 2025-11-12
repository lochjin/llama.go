package routes

import (
	"context"
	"github.com/Qitmeer/llama.go/api"
)

func PullModel(ctx context.Context, name string, fn func(api.ProgressResponse)) error {
	fn(api.ProgressResponse{Status: "success"})
	return nil
}
