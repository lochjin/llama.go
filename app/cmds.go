package app

import (
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"os"
)

func commands() []*cli.Command {
	cmds := []*cli.Command{}
	cmds = append(cmds, downloadCmd())
	cmds = append(cmds, embeddingCmd())
	return cmds
}

func downloadCmd() *cli.Command {
	return &cli.Command{
		Name:        "download",
		Aliases:     []string{"r"},
		Category:    "llama",
		Usage:       "Download model",
		Description: "Download model",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}

func embeddingCmd() *cli.Command {
	return &cli.Command{
		Name:        "embedding",
		Aliases:     []string{"e"},
		Category:    "llama",
		Usage:       "Generate high-dimensional embedding vector of a given text",
		Description: "Generate high-dimensional embedding vector of a given text",
		Action: func(ctx *cli.Context) error {
			cfg := config.Conf
			err := initLog(cfg)
			if err != nil {
				return err
			}
			log.Info("Start embedding")
			err = cfg.Load()
			if err != nil {
				return err
			}
			ret, err := wrapper.LlamaEmbedding(cfg, cfg.ModelPath(), cfg.Prompt, cfg.EmbdOutputFormat)
			if err != nil {
				return err
			}
			if len(cfg.OutputFile) > 0 {
				return saveOutputToFile(cfg.OutputFile, ret)
			} else {
				fmt.Println("result:")
				fmt.Println(ret)
			}
			return nil
		},
	}
}

func saveOutputToFile(outFilePath string, content string) error {
	outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		outFile.Close()
	}()
	_, err = outFile.WriteString(content)
	return err
}
