// Copyright (c) 2017-2025 The qitmeer developers

package config

import (
	"fmt"
	"github.com/Qitmeer/llama.go/common"
	"github.com/Qitmeer/llama.go/model"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"math"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultLogLevel = "info"
	defaultNPredict = -1
	DefaultHost     = "127.0.0.1:8081"
	DefaultPort     = "8081"
	DefaultModelDir = "./data/models"
)

var (
	defaultHomeDir     = "."
	defaultSwaggerFile = filepath.Join(defaultHomeDir, "swagger.json")
)

var (
	Conf = &Config{}

	LogLevel = &cli.StringFlag{
		Name:        "log-level",
		Aliases:     []string{"l"},
		Usage:       "Logging level {trace, debug, info, warn, error}",
		Value:       defaultLogLevel,
		Destination: &Conf.LogLevel,
	}

	Model = &cli.StringFlag{
		Name:        "model",
		Aliases:     []string{"m"},
		Usage:       "The name of the model file located in the 'model-dir' repository path or absolute path of resource file",
		Destination: &Conf.Model,
	}

	ModelDir = &cli.StringFlag{
		Name:        "model-dir",
		Aliases:     []string{"md"},
		Usage:       "Path for storing model files",
		Value:       DefaultModelDir,
		Destination: &Conf.ModelDir,
		EnvVars:     []string{"LLAMAGO_MODEL_DIR"},
	}

	CtxSize = &cli.IntFlag{
		Name:        "ctx-size",
		Aliases:     []string{"c"},
		Usage:       "Set the size of the prompt context. The default is 4096, but if a LLaMA model was built with a longer context, increasing this value will provide better results for longer input/inference",
		Value:       4096,
		Destination: &Conf.CtxSize,
	}

	Prompt = &cli.StringFlag{
		Name:        "prompt",
		Aliases:     []string{"p"},
		Usage:       "Provide a prompt directly as a command-line option.",
		Destination: &Conf.Prompt,
	}

	NGpuLayers = &cli.IntFlag{
		Name:        "n-gpu-layers",
		Aliases:     []string{"ngl"},
		Usage:       "When compiled with GPU support, this option allows offloading some layers to the GPU for computation. Generally results in increased performance.",
		Value:       -1,
		Destination: &Conf.NGpuLayers,
	}

	NPredict = &cli.IntFlag{
		Name:        "n-predict",
		Aliases:     []string{"n"},
		Usage:       "Set the number of tokens to predict when generating text. Adjusting this value can influence the length of the generated text.",
		Value:       defaultNPredict,
		Destination: &Conf.NPredict,
	}

	Interactive = &cli.BoolFlag{
		Name:        "interactive",
		Aliases:     []string{"i"},
		Usage:       "Run the program in interactive mode, allowing you to provide input directly and receive real-time responses",
		Value:       false,
		Destination: &Conf.Interactive,
	}

	Seed = &cli.UintFlag{
		Name:        "seed",
		Aliases:     []string{"s"},
		Usage:       "Set the random number generator (RNG) seed (default: -1, -1 = random seed).",
		Value:       math.MaxUint32,
		Destination: &Conf.Seed,
	}

	Pooling = &cli.StringFlag{
		Name:        "pooling",
		Aliases:     []string{"o"},
		Usage:       "pooling type for embeddings, use model default if unspecified {none,mean,cls,last,rank}",
		Value:       "mean",
		Destination: &Conf.Pooling,
	}

	EmbdNormalize = &cli.IntFlag{
		Name:        "embd-normalize",
		Aliases:     []string{"N"},
		Usage:       "normalisation for embeddings (default: %d) (-1=none, 0=max absolute int16, 1=taxicab, 2=euclidean, >2=p-norm)",
		Value:       2,
		Destination: &Conf.EmbdNormalize,
	}

	EmbdOutputFormat = &cli.StringFlag{
		Name:        "embd-output-format",
		Aliases:     []string{"FORMAT"},
		Usage:       "empty = default, \"array\" = [[],[]...], \"json\" = openai style, \"json+\" = same \"json\" + cosine similarity matrix",
		Value:       "json",
		Destination: &Conf.EmbdOutputFormat,
	}

	EmbdSeparator = &cli.StringFlag{
		Name:        "embd-separator",
		Aliases:     []string{"STRING"},
		Usage:       "separator of embeddings (default \\n) for example \"<#sep#>\\",
		Value:       "<#sep#>",
		Destination: &Conf.EmbdSeparator,
	}

	BatchSize = &cli.IntFlag{
		Name:        "batch-size",
		Aliases:     []string{"b"},
		Usage:       "logical maximum batch size",
		Value:       2048,
		Destination: &Conf.BatchSize,
	}

	UBatchSize = &cli.IntFlag{
		Name:        "ubatch-size",
		Aliases:     []string{"ub"},
		Usage:       "physical maximum batch size",
		Value:       512,
		Destination: &Conf.UBatchSize,
	}

	OutputFile = &cli.StringFlag{
		Name:        "output-file",
		Aliases:     []string{"of"},
		Usage:       "output file",
		Destination: &Conf.OutputFile,
	}

	Host = &cli.StringFlag{
		Name:        "host",
		Aliases:     []string{"ho"},
		Usage:       fmt.Sprintf("IP Address for the ollama server (default %s)", DefaultHost),
		Value:       DefaultHost,
		EnvVars:     []string{"LLAMAGO_HOST"},
		Destination: &Conf.Host,
	}

	Origins = &cli.StringFlag{
		Name:        "origins",
		Aliases:     []string{"or"},
		Usage:       "A comma separated list of allowed origins",
		EnvVars:     []string{"LLAMAGO_ORIGINS"},
		Destination: &Conf.Origins,
	}

	AppFlags = []cli.Flag{
		LogLevel,
		Model,
		ModelDir,
		CtxSize,
		Prompt,
		NGpuLayers,
		NPredict,
		Interactive,
		Seed,
		Pooling,
		EmbdNormalize,
		EmbdOutputFormat,
		EmbdSeparator,
		BatchSize,
		UBatchSize,
		OutputFile,
		Host,
		Origins,
	}
)

type Config struct {
	LogLevel string
	Model    string
	ModelDir string

	CtxSize          int
	Prompt           string
	NGpuLayers       int
	NPredict         int
	Interactive      bool
	Seed             uint
	Pooling          string
	EmbdNormalize    int
	EmbdOutputFormat string
	EmbdSeparator    string
	BatchSize        int
	UBatchSize       int
	OutputFile       string
	Host             string
	Origins          string
}

func (c *Config) Load() error {
	log.Debug("Try to load config")
	if !c.HasModel() {
		return fmt.Errorf("No config model")
	}
	log.Debug("Model info", "model path", c.ModelPath())
	return nil
}

func (c *Config) ModelPath() string {
	if len(c.Model) <= 0 {
		return ""
	}
	if !strings.Contains(c.Model, model.EXT) {
		return ""
	}
	if common.IsFilePath(c.Model) {
		return c.Model
	}
	ret, err := filepath.Abs(filepath.Join(c.ModelDir, c.Model))
	if err != nil {
		return ""
	}
	return ret
}

func (c *Config) HasModel() bool {
	return len(c.ModelPath()) > 0
}

func (c *Config) GetModelFileInfos() []os.FileInfo {
	var firstInfo os.FileInfo
	if c.HasModel() {
		info, err := os.Stat(c.ModelPath())
		if err == nil {
			firstInfo = info
		}
	}
	hasAdd := false
	ret := []os.FileInfo{}

	if common.IsExist(c.ModelDir) {
		entries, err := os.ReadDir(c.ModelDir)
		if err != nil {
			log.Error(err.Error())
		} else {
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
				ret = append(ret, info)
				if firstInfo != nil && !hasAdd {
					if os.SameFile(firstInfo, info) {
						hasAdd = true
					}
				}
			}
		}
	}

	if firstInfo != nil && !hasAdd {
		ret = append(ret, firstInfo)
	}
	return ret
}

func (c *Config) IsLonely() bool {
	return len(c.Prompt) > 0 || c.Interactive
}

func (c *Config) HostURL() *url.URL {
	defaultPort := DefaultPort
	chost := c.Host
	scheme, hostport, ok := strings.Cut(chost, "://")
	switch {
	case !ok:
		scheme, hostport = "http", chost
	case scheme == "http":
		defaultPort = "80"
	case scheme == "https":
		defaultPort = "443"
	}

	hostport, path, _ := strings.Cut(hostport, "/")
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host, port = "127.0.0.1", defaultPort
		if ip := net.ParseIP(strings.Trim(hostport, "[]")); ip != nil {
			host = ip.String()
		} else if hostport != "" {
			host = hostport
		}
	}

	if n, err := strconv.ParseInt(port, 10, 32); err != nil || n > 65535 || n < 0 {
		log.Warn("invalid port, using default", "port", port, "default", defaultPort)
		port = defaultPort
	}

	return &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, port),
		Path:   path,
	}
}

// AllowedOrigins returns a list of allowed origins. AllowedOrigins can be configured via the LLAMAGO_ORIGINS environment variable.
func (c *Config) AllowedOrigins() []string {
	origins := []string{}
	if len(c.Origins) > 0 {
		origins = strings.Split(c.Origins, ",")
	}

	for _, origin := range []string{"localhost", "127.0.0.1", "0.0.0.0"} {
		origins = append(origins,
			fmt.Sprintf("http://%s", origin),
			fmt.Sprintf("https://%s", origin),
			fmt.Sprintf("http://%s", net.JoinHostPort(origin, "*")),
			fmt.Sprintf("https://%s", net.JoinHostPort(origin, "*")),
		)
	}

	origins = append(origins,
		"app://*",
		"file://*",
		"tauri://*",
		"vscode-webview://*",
		"vscode-file://*",
	)

	return origins
}
