#include "process.h"
#include "log.h"
#include "whisper_service.h"
#include "server/server.h"

#include <sstream>
#include <string>

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

Result llama_gen(int id, const char * js_str) {
    if (!Server::instance().is_running()) {
        return {false, nullptr};
    }
    if (!js_str) {
        return {false, nullptr};
    }

    server_http_req rq{
            id,
            std::string(js_str),
            [](int cid, const std::string & content) {
                PushToChan(cid, content.c_str());
                return true;
            }
    };

    server_http_res_ptr rp = Server::instance().post_completions(rq);
    Server::flush_http_response_to_sink(rq, *rp);
    const bool ok = rp->is_success();

    CloseChan(id);
    return {ok, nullptr};
}

Result llama_chat(int id, const char * js_str) {
    (void)id;
    (void)js_str;
    return {true, nullptr};
}

Result whisper_gen(const char * model,const char * input) {
    (void)model;
    (void)input;
    return {true, nullptr};
}

CommonParams get_common_params() {
    return {};
}

Result get_props() {
    return {true, nullptr};
}

Result get_slots() {
    return {true, nullptr};
}