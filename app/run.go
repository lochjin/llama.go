// Copyright (c) 2017-2025 The qitmeer developers

package app

import (
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/version"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"os"
)

func Run() error {
	app := &cli.App{
		Name:    "",
		Version: version.String(),
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Qitmeer",
			},
		},
		Copyright:            "(c) 2025 Qitmeer",
		Usage:                "Llama",
		Flags:                config.AppFlags,
		EnableBashCompletion: true,
		Commands:             commands(),
		Action: func(c *cli.Context) error {
			print(version.String())
			return nil
		},
	}

	return app.Run(os.Args)
}

func OnBefore(ctx *cli.Context) error {
	err := initLog(config.Conf)
	if err != nil {
		return err
	}
	log.Info("Before init")
	err = config.Conf.Load()
	if err != nil {
		return err
	}
	return nil
}
