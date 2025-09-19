package app

import (
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/server"
	"github.com/Qitmeer/llama.go/system"
	"github.com/Qitmeer/llama.go/system/limits"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
)

func commands() []*cli.Command {
	cmds := []*cli.Command{}
	cmds = append(cmds, serveCmd())
	cmds = append(cmds, runCmd())
	cmds = append(cmds, downloadCmd())
	cmds = append(cmds, embeddingCmd())
	cmds = append(cmds, whisperCmd())
	return cmds
}

func serveCmd() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Aliases:     []string{"s"},
		Category:    "llama",
		Usage:       "llama.go server",
		Description: "llama.go server",
		Before:      OnBefore,
		Action: func(ctx *cli.Context) error {
			err := limits.SetLimits()
			if err != nil {
				return err
			}
			interrupt := system.InterruptListener()
			cfg := config.Conf

			ser := server.New(ctx, cfg)

			err = wrapper.LlamaStart(cfg)
			if err != nil {
				log.Error(err.Error())
			}
			log.Info("Started llama core")

			err = ser.Start()
			defer func() {
				err = ser.Stop()
				if err != nil {
					log.Error(err.Error())
				}
				err = wrapper.LlamaStop()
				if err != nil {
					log.Error(err.Error())
				}
			}()

			if err != nil {
				return err
			}
			<-interrupt

			return nil
		},
	}
}

func runCmd() *cli.Command {
	return &cli.Command{
		Name:        "run",
		Aliases:     []string{"r"},
		Category:    "llama",
		Usage:       "llama.go run",
		Description: "llama.go run",
		Before:      OnBefore,
		Action: func(ctx *cli.Context) error {
			err := limits.SetLimits()
			if err != nil {
				return err
			}
			interrupt := system.InterruptListener()
			cfg := config.Conf

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				err = wrapper.LlamaStartInteractive(cfg)
				if err != nil {
					log.Error(err.Error())
				}
			}()
			<-interrupt

			log.Info("Stop run cmd")
			err = wrapper.LlamaStopInteractive()
			if err != nil {
				log.Error(err.Error())
			}
			wg.Wait()
			log.Info("Stopped run cmd")
			return nil
		},
	}
}

func downloadCmd() *cli.Command {
	return &cli.Command{
		Name:        "download",
		Aliases:     []string{"d"},
		Category:    "llama",
		Usage:       "Download model",
		Description: "Download model",
		Before:      OnBefore,
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
		Before:      OnBefore,
		Action: func(ctx *cli.Context) error {
			cfg := config.Conf
			log.Info("Start embedding")
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

func whisperCmd() *cli.Command {
	return &cli.Command{
		Name:        "whisper",
		Aliases:     []string{"w"},
		Category:    "whisper",
		Usage:       "Generate text by whisper model",
		Description: "Generate text by whisper model",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input file path for generate.",
			},
		},
		Before: OnBefore,
		Action: func(ctx *cli.Context) error {
			cfg := config.Conf
			log.Info("Start whisper")
			if !ctx.IsSet("input") {
				return fmt.Errorf("No input file")
			}
			ret, err := wrapper.WhisperGenerate(cfg, ctx.Value("input").(string))
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
