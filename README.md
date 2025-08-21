# llama.go
Go bindings to llama.cpp

### Installation
***make sure you have `git golang cmake gcc make` installed on the system to build.***
* Build from source
```bash
~ git clone https://github.com/Qitmeer/llama.go.git
~ cd llama.go
~ make
```

### Get model
* Manually download the model:[Hugging Face Qwen3-8B-GGUF](https://huggingface.co/ggml-org/Qwen3-8B-GGUF/tree/main)
* Please first set the storage location of the model file, which can be done using environment variables `LLAMAGO_MODEL_DIR` or command-line parameters `model-dir`
* Default model files directory is `./data/models`

```bash
~ ./llama --model-dir=<your_model_files_directory>
or
~ export LLAMAGO_MODEL_DIR=<your_model_files_directory>
```

### Local startup

```bash
~ ./llama --model=qwen2.5-0.5b-q8_0.gguf --prompt=天空为什么是蓝的
```
Or enable interactive mode to run:
```bash
~ ./llama --model=qwen2.5-0.5b-q8_0.gguf -i
```


### As the startup of the server

```bash
~ ./llama --model=qwen2.5-0.5b-q8_0.gguf
```

* Support REST API:
```bash
~ curl -s -k -X POST -H 'Content-Type: application/json' --data '{"prompt":"天空为什么是蓝的"}' http://127.0.0.1:8081/api/generate
```

#### WebUI
* Enter this address `http://127.0.0.1:8081` in the browser

### Embedding

* Local mode:
```bash
~ ./llama --model=qwen2.5-0.5b-q8_0.gguf --prompt=天空为什么是蓝的 --output-file=./embs.json embedding
```

* Server mode:
```bash
~ curl -s -k -X POST -H 'Content-Type: application/json' --data '{"input":["天空","蓝色"]}' http://127.0.0.1:8081/api/embed
~ curl -s -k -X POST -H 'Content-Type: application/json' --data '{"prompt":"天空为什么是蓝的"}' http://127.0.0.1:8081/api/embeddings
```

