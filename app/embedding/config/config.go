// Copyright (c) 2017-2025 The qitmeer developers

package config

import (
	"github.com/urfave/cli/v2"
)

const (
	DefaultEmbdNormalize    = 2
	DefaultEmbdOutputFormat = "json"
	DefaultEmbdSeparator    = "<#sep#>"
)

var (
	Conf = &Config{}

	EmbdNormalize = &cli.IntFlag{
		Name:        "embd-normalize",
		Aliases:     []string{"N"},
		Usage:       "normalisation for embeddings (default: %d) (-1=none, 0=max absolute int16, 1=taxicab, 2=euclidean, >2=p-norm)",
		Value:       DefaultEmbdNormalize,
		Destination: &Conf.EmbdNormalize,
	}

	EmbdOutputFormat = &cli.StringFlag{
		Name:        "embd-output-format",
		Aliases:     []string{"FORMAT"},
		Usage:       "empty = default, \"array\" = [[],[]...], \"json\" = openai style, \"json+\" = same \"json\" + cosine similarity matrix",
		Value:       DefaultEmbdOutputFormat,
		Destination: &Conf.EmbdOutputFormat,
	}

	EmbdSeparator = &cli.StringFlag{
		Name:        "embd-separator",
		Aliases:     []string{"STRING"},
		Usage:       "separator of embeddings (default \\n) for example \"<#sep#>\\",
		Value:       DefaultEmbdSeparator,
		Destination: &Conf.EmbdSeparator,
	}

	AppFlags = []cli.Flag{
		EmbdNormalize,
		EmbdOutputFormat,
		EmbdSeparator,
	}
)

type Config struct {
	EmbdNormalize    int
	EmbdOutputFormat string
	EmbdSeparator    string
}
