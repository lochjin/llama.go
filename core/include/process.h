#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

typedef struct Result {
    bool ret;
    const char *content;
} Result;

typedef struct CommonParams {
    bool endpoint_props;
}CommonParams;

bool llama_start(const char * args);
bool llama_stop();
Result llama_gen(int id,const char * model,const char * js_str);
Result llama_chat(int id,const char * js_str);

bool llama_interactive_start(const char * args,const char * prompt);
bool llama_interactive_stop();

Result whisper_gen(const char * model,const char * input);

CommonParams get_common_params();
Result get_props();
Result get_slots();

#ifdef __cplusplus
}
#endif