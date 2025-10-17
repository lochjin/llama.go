// Copyright (c) 2017-2025 The qitmeer developers

package run

import (
	"github.com/urfave/cli/v2"
)

var (
	Conf = &Config{}

	Keepalive = &cli.StringFlag{
		Name:        "keepalive",
		Aliases:     []string{"k"},
		Usage:       "Duration to keep a model loaded (e.g. 5m)",
		Destination: &Conf.Keepalive,
	}

	Verbose = &cli.BoolFlag{
		Name:        "verbose",
		Aliases:     []string{"v"},
		Usage:       "Show timings for response",
		Value:       false,
		Destination: &Conf.Verbose,
	}

	Insecure = &cli.BoolFlag{
		Name:        "insecure",
		Aliases:     []string{"i"},
		Usage:       "Use an insecure registry",
		Value:       false,
		Destination: &Conf.Insecure,
	}

	Nowordwrap = &cli.BoolFlag{
		Name:        "nowordwrap",
		Aliases:     []string{"n"},
		Usage:       "Don't wrap words to the next line automatically",
		Value:       false,
		Destination: &Conf.Nowordwrap,
	}

	Hidethinking = &cli.BoolFlag{
		Name:        "hidethinking",
		Aliases:     []string{"ht"},
		Usage:       "Hide thinking output (if provided)",
		Value:       false,
		Destination: &Conf.Hidethinking,
	}

	Format = &cli.StringFlag{
		Name:        "format",
		Aliases:     []string{"f"},
		Usage:       "Response format (e.g. json)",
		Destination: &Conf.Format,
	}

	Think = &cli.StringFlag{
		Name:        "think",
		Aliases:     []string{"t"},
		Usage:       "Enable thinking mode: true/false or high/medium/low for supported models",
		Destination: &Conf.Think,
	}

	Prompt = &cli.StringFlag{
		Name:        "prompt",
		Aliases:     []string{"p"},
		Usage:       "Provide a prompt directly as a command-line option.",
		Destination: &Conf.Prompt,
	}

	AppFlags = []cli.Flag{
		Keepalive,
		Verbose,
		Insecure,
		Nowordwrap,
		Hidethinking,
		Format,
		Think,
		Prompt,
	}
)

type Config struct {
	Keepalive    string
	Verbose      bool
	Insecure     bool
	Nowordwrap   bool
	Hidethinking bool
	Format       string
	Think        string
	Prompt       string
}
