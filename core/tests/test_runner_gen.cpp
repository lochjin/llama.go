#include <iostream>
#include <thread>
#include <chrono>
#include <cstdlib>
#include <future>

#include "./../src/scheduler.h"
#include "log.h"

int main() {
    common_log_verbosity_thold=1;
    const char* env_model = "LLAMA_TEST_MODEL";
    const char* model = std::getenv(env_model);

    if (model == nullptr) {
        std::cerr << "error：can't find " << env_model << std::endl;
        return EXIT_FAILURE;
    }

    std::cout << "env: " << env_model << "=" << model << std::endl;

    std::vector<std::string> v_args;
    v_args.push_back("test_runner_gen");
    v_args.push_back("-m");
    v_args.push_back(model);
    v_args.push_back("--seed");
    v_args.push_back("0");

    if (!Scheduler::instance().start(v_args)) {
        return EXIT_FAILURE;
    }

    std::string js_str="{\"prompt\":\"why the sky is blue\"}";
    //std::string js_str="{\"prompt\":\"why the sky is blue\",\"stream\":true}";

    int id=1;
    Request rq{id,std::string(js_str)};
    Response rp{id};

    rp.write = [](int id, const std::string& content) {
        std::cout<<content;
        return true;
    };
    rp.is_writable = [](int id) {
        return true;
    };
    rp.complete = [&](int id) {};

    Scheduler::instance().handle_completions_oai(rq,rp);
    if (!rp.success) {
        return EXIT_FAILURE;
    }

    std::cout<<"stop:"<<Scheduler::instance().stop()<<std::endl;
    return EXIT_SUCCESS;
}