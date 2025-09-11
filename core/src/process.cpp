#include "process.h"
#include "runner.h"
#include "log.h"
#include "whisper_service.h"
#include "scheduler.h"

extern "C" {
    void PushToChan(int id, const char* val);
    void CloseChan(int id);
}

int llama_start(const char * args) {
    if (Scheduler::instance().is_running()) {
        return EXIT_FAILURE;
    }

    std::istringstream iss(args);
    std::vector<std::string> v_args;
    std::string v_a;
    while (iss >> v_a) {
        v_args.push_back(v_a);
    }

    if (Scheduler::instance().start(v_args)) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}

int llama_stop() {
    if (!Scheduler::instance().is_running()) {
        return EXIT_FAILURE;
    }
    if (Scheduler::instance().stop()) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}

Result llama_gen(const char * js_str) {
    if (!Scheduler::instance().is_running()) {
        return {false,""};
    }
    Request rq{std::string(js_str)};
    Response rp;
    Scheduler::instance().handle_completions_oai(rq,rp);
    if (!rp.success) {
        return {false,""};
    }

    char* arr = new char[rp.content.size() + 1];
    std::copy(rp.content.begin(), rp.content.end(), arr);
    arr[rp.content.size()] = '\0';

    return {true,arr};
}

Result llama_chat(const char * js_str) {
    if (!Scheduler::instance().is_running()) {
        return {false,""};
    }
    return {false,""};
}

Result whisper_gen(const char * model,const char * input) {
    WhisperService ws;

    std::string result = ws.generate(std::string(model),std::string(input));
    if (result.empty()) {
        return {false,""};
    }
    char* arr = new char[result.size() + 1];
    std::copy(result.begin(), result.end(), arr);
    arr[result.size()] = '\0';

    return {true,arr};
}