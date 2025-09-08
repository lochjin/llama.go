package wrapper

/*
#include <stdlib.h>
#include "core.h"
*/
import "C"

import (
	"fmt"
	"github.com/ollama/ollama/api"
	"sync"
	"unsafe"

	"github.com/Qitmeer/llama.go/config"
)

var (
	mu         sync.Mutex
	channels   = make(map[int]chan any)
	nextChanID = 0
)

func LlamaInteractive(cfg *config.Config) error {
	if !cfg.Interactive {
		return fmt.Errorf("Not config interactive")
	}
	if !cfg.HasModel() {
		return fmt.Errorf("No model")
	}
	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	cfgArgs := fmt.Sprintf("llama -i --model %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d",
		cfg.ModelPath(), cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed)
	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ret := C.llama_start(ca, 0, ip)
	if ret != 0 {
		return fmt.Errorf("Llama start error")
	}
	ret = C.llama_stop()
	if ret != 0 {
		return fmt.Errorf("Llama stop error")
	}
	return nil
}

func LlamaGenerate(prompt string) (string, error) {
	if len(prompt) <= 0 {
		return "", fmt.Errorf("No prompt")
	}
	ip := C.CString(prompt)
	defer C.free(unsafe.Pointer(ip))

	ret := C.llama_gen(ip)
	if ret == nil {
		return "", fmt.Errorf("Llama run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}

func LlamaChat(msgs []api.Message) (string, error) {
	size := len(msgs)
	if size <= 0 {
		return "", fmt.Errorf("No messages for chat")
	}
	roles := make([]*C.char, size)
	contents := make([]*C.char, size)

	for i, m := range msgs {
		roles[i] = C.CString(m.Role)
		defer C.free(unsafe.Pointer(roles[i]))

		contents[i] = C.CString(m.Content)
		defer C.free(unsafe.Pointer(contents[i]))
	}

	rolesPtr := (**C.char)(unsafe.Pointer(&roles[0]))
	contentsPtr := (**C.char)(unsafe.Pointer(&contents[0]))

	ret := C.llama_chat(rolesPtr, contentsPtr, C.int(size))
	if ret == nil {
		return "", fmt.Errorf("Llama run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}

func LlamaStart(cfg *config.Config) error {
	if !cfg.HasModel() {
		return fmt.Errorf("No model")
	}
	cfgArgs := fmt.Sprintf("llama -i --model %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d",
		cfg.ModelPath(), cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed)
	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	ret := C.llama_start(ca, 1, ip)
	if ret != 0 {
		return fmt.Errorf("Llama start error")
	}
	return nil
}

func LlamaStop() error {
	ret := C.llama_stop()
	if ret != 0 {
		return fmt.Errorf("Llama stop error")
	}
	return nil
}

func LlamaEmbedding(cfg *config.Config, model string, prompts string, embdOutputFormat string) (string, error) {
	if len(model) <= 0 {
		return "", fmt.Errorf("No model")
	}
	if len(prompts) <= 0 {
		return "", fmt.Errorf("No prompt")
	}
	ip := C.CString(prompts)
	defer C.free(unsafe.Pointer(ip))

	cfgArgs := fmt.Sprintf("llama --model %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d --embd-normalize %d --batch-size %d --ubatch-size %d",
		model, cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed, cfg.EmbdNormalize, cfg.BatchSize, cfg.UBatchSize)
	if len(cfg.Pooling) > 0 {
		cfgArgs = fmt.Sprintf("%s --pooling %s", cfgArgs, cfg.Pooling)
	}
	if len(embdOutputFormat) > 0 {
		cfgArgs = fmt.Sprintf("%s --embd-output-format %s", cfgArgs, embdOutputFormat)
	}
	if len(cfg.EmbdSeparator) > 0 {
		cfgArgs = fmt.Sprintf("%s --embd-separator %s", cfgArgs, cfg.EmbdSeparator)
	}
	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ret := C.llama_embedding(ca, ip)
	if ret == nil {
		return "", fmt.Errorf("llama_embedding run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}

func WhisperGenerate(cfg *config.Config, input string) (string, error) {
	if !cfg.HasModel() {
		return "", fmt.Errorf("No model")
	}
	if len(input) <= 0 {
		return "", fmt.Errorf("No input")
	}
	model := C.CString(cfg.ModelPath())
	defer C.free(unsafe.Pointer(model))

	ip := C.CString(input)
	defer C.free(unsafe.Pointer(ip))

	ret := C.whisper_gen(model, ip)
	if ret == nil {
		return "", fmt.Errorf("Whisper run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}

//export PushToChan
func PushToChan(id C.int, val *C.char) {
	str := C.GoString(val)
	mu.Lock()
	ch, ok := channels[int(id)]
	mu.Unlock()
	if ok {
		ch <- str
	}
}

//export CloseChan
func CloseChan(id C.int) {
	mu.Lock()
	ch, ok := channels[int(id)]
	if ok {
		close(ch)
		delete(channels, int(id))
	}
	mu.Unlock()
}
