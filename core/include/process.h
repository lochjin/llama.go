#pragma once

#ifdef __cplusplus
extern "C" {
#endif

    struct Result {
        const char *content;
        int ret;
    };

    int llama_start(const char * args);
    int llama_stop();
    const char * llama_gen(const char * js_str);
    const char * llama_chat(const char * js_str);
    const char * whisper_gen(const char * model,const char * input);

#ifdef __cplusplus
}
#endif