package runner

import (
	"errors"
	"fmt"

	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

type Service struct {
	ctx     *cli.Context
	cfg     *config.Config
	running bool
}

func New(ctx *cli.Context, cfg *config.Config) *Service {
	log.Info("New Runner ...")
	ser := Service{ctx: ctx, cfg: cfg, running: false}
	return &ser
}

func (s *Service) Start() error {
	if s.IsRunning() {
		return errors.New("Already Running")
	}
	log.Info("Start Runner...")
	err := wrapper.LlamaStart(s.cfg)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Started llama core")
	s.running = true
	return nil
}

func (s *Service) Stop() error {
	if !s.IsRunning() {
		return errors.New("Not running")
	}
	log.Info("Stop Runner...")
	err := wrapper.LlamaStop()
	if err != nil {
		log.Error(err.Error())
	}
	s.running = false
	return nil
}

func (s *Service) IsRunning() bool {
	return s.running
}

func (s *Service) Generate(id int, prompt string, stream bool) error {
	return wrapper.LlamaGenerate(id, fmt.Sprintf("{\"prompt\":\"%s\",\"stream\":%v}", prompt, stream))
}

func (s *Service) Chat(id int, jsStr string) error {
	return wrapper.LlamaChat(id, jsStr)
}
