#include "process.h"
#include "log.h"
#include "whisper_service.h"
#include "server/server.h"

extern "C" {
    void PushToChan(int id, const char* val);
    void CloseChan(int id);
}

bool llama_start(const char * args) {
    if (Server::instance().is_running()) {
        return false;
    }

    std::istringstream iss(args);
    std::vector<std::string> v_args;
    std::string v_a;
    while (iss >> v_a) {
        v_args.push_back(v_a);
    }

    if (!Server::instance().start(v_args)) {
        return false;
    }
    return true;
}

bool llama_stop() {
    if (!Server::instance().is_running()) {
        return false;
    }
    if (!Server::instance().stop()) {
        return false;
    }
    return true;
}

Result llama_gen(int id,const char * js_str) {
    if (!Server::instance().is_running()) {
        return {false};
    }
    server_http_req rq{
            id,
            std::string(js_str),
            [](int id) {CloseChan(id);},
            [](int id, const std::string& content) {
                PushToChan(id, content.c_str());
                return true;
            }
    };

    server_http_res_ptr rp=Server::instance().post_completions(rq);

    return {rp->is_success()};
}

Result llama_chat(int id,const const char * js_str) {


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