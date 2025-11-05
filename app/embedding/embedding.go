package embedding

import (
	"fmt"

	"github.com/Qitmeer/llama.go/common"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

func EmbeddingHandler(ctx *cli.Context) error {
	var prompts string
	if ctx.Args().Len() > 0 {
		prompts = ctx.Args().Slice()[0]
	}
	if len(prompts) <= 0 {
		return fmt.Errorf("No prompt")
	}
	cfg := config.Conf
	log.Info("Start embedding")
	ret, err := wrapper.LlamaEmbedding(cfg, prompts, "")
	if err != nil {
		return err
	}
	if len(cfg.OutputFile) > 0 {
		return common.SaveOutputToFile(cfg.OutputFile, ret)
	} else {
		fmt.Println("result:")
		fmt.Println(ret)
	}
	return nil
}
