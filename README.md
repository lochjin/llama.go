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

### As the startup of the server

```bash
~ ./llama --model=qwen2.5-0.5b-q8_0.gguf serve
or
~ ./llama --model=gpt-oss-20b-mxfp4.gguf --jinja serve
```

### client:

```bash
~ ./llama run 天空为什么是蓝的
```
Or enable interactive mode to run:
```bash
~ ./llama run
```
 
#### Download Model by CLI:
```bash
~ ./llama pull gte-small-Q8_0-GGUF
```
or
```bash
~ ./llama pull gte-small-Q8_0-GGUF:gte-small-q8_0.gguf
```
or
```bash
~ ./llama pull llamago/gte-small-Q8_0-GGUF:gte-small-q8_0.gguf
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
~ ./llama --model=qwen2.5-0.5b-q8_0.gguf embedding 天空为什么是蓝的 --output-file=./embs.json
```

* Server mode:
```bash
~ curl -s -k -X POST -H 'Content-Type: application/json' --data '{"input":["天空","蓝色"]}' http://127.0.0.1:8081/api/embed
~ curl -s -k -X POST -H 'Content-Type: application/json' --data '{"prompt":"天空为什么是蓝的"}' http://127.0.0.1:8081/api/embeddings
```
### Whisper
* Firstly, you need to download the model from this address `https://huggingface.co/ggerganov/whisper.cpp` and then place it in `LLAMAGO_MODEL_DIR` or `model-dir`

```bash
~ ./llama --model=ggml-base.en.bin whisper --input="./your-voice.wav"
```



