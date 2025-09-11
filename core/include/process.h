#pragma once

#ifdef __cplusplus
extern "C" {
#endif

    #include <stdbool.h>

    typedef struct Result {
        bool ret;
        const char *content;
    } Result;

    int llama_start(const char * args);
    int llama_stop();
    Result llama_gen(const char * js_str);
    Result llama_chat(const char * js_str);
    Result whisper_gen(const char * model,const char * input);

#ifdef __cplusplus
}
#endif