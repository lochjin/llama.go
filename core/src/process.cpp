#include "process.h"
#include "runner.h"
#include "log.h"
#include "whisper_service.h"

extern "C" {
    void PushToChan(int id, const char* val);
    void CloseChan(int id);
}

bool llama_start(const char * args) {
    return true;
}

bool llama_stop() {

    return true;
}

Result llama_gen(int id,const char * model,const char * js_str) {

    return {true};
}

Result llama_chat(int id,const char * model,const char * js_str) {


    return {true};
}

Result whisper_gen(const char * model,const char * input) {
    return {true};
}

CommonParams get_common_params() {
    return {};
}

Result get_props() {
    return {true};
}

Result get_slots() {

    return {true};
}