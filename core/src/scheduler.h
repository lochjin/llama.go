#pragma once

#include <map>
#include <memory>
#include <mutex>
#include <thread>
#include <vector>

#include "server_context.h"
#include "singleton.h"

struct Request {
    int id;
    std::string model;
    std::string body;
    std::function<bool()> is_connection_closed = []() { return false; };
};

struct Response {
    int id;
    bool success;
    std::function<void(int)> complete;
    std::function<bool(int,const std::string&)> write;
    std::function<bool(int)> is_writable;
};

class Scheduler : public patterns::Singleton<Scheduler> {
    friend class patterns::Singleton<Scheduler>;

private:
    struct CtxEntry {
        server_context ctx;
        std::thread    thread;
    };

    // model_key -> context/worker-thread
    std::map<std::string, std::unique_ptr<CtxEntry>> ctx_servers;
    std::mutex ctx_mu;
    std::vector<std::string> base_args;
    std::string default_model_key;
    bool running= false;
    Scheduler();
    ~Scheduler();

    server_context * get_or_create_ctx(const std::string & model_key);
    server_context * get_default_ctx();
    bool init_server_context_for_model(const std::string & model_key);

public:
    bool start(const std::vector<std::string>& args);
    bool stop();
    void cleanup();

    bool is_running();

    bool init_server_context(const std::vector<std::string>& args);
    void handle_completions(const Request & req, Response & res);
    void handle_completions_impl(server_task_type type,std::string model,json & data,const std::vector<raw_buffer> & files,const std::function<bool()> & is_connection_closed,Response & res,oaicompat_type oaicompat);
    void handle_completions_oai(const Request & req, Response & res);
    void handle_chat_completions(const Request & req, Response & res);

    void handle_embeddings_impl(const Request & req, Response & res, oaicompat_type oaicompat);
    void handle_embeddings(const Request & req, Response & res);
    void handle_embeddings_oai(const Request & req, Response & res);
    void res_error(Response & res, const json & error_data);
    void res_ok(Response & res, const json & data);
    common_params *get_common_params();
    std::string get_props();
    std::string get_slots(bool fail_on_no_slot= false);
};