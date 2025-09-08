#include "process.h"
#include "runner.h"
#include "log.h"
#include "whisper_service.h"
#include "scheduler.h"

static Runner *g_runner;
static int g_idx=0;

int llama_start(const char * args,int async,const char * prompt) {
    if (g_runner != nullptr) {
        LOG("Delete last runner: id=%d\n",g_runner->getID());
        delete g_runner;
        g_runner= nullptr;
    }
    std::istringstream iss(args);
    std::vector<std::string> v_args;
    std::string v_a;
    while (iss >> v_a) {
        v_args.push_back(v_a);
    }

    g_runner=new Runner(g_idx,v_args,async>0,std::string(prompt));
    g_idx++;
    if (g_runner->start()) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}

int llama_stop() {
    if (g_runner == nullptr) {
        LOG("Runner is already delete\n");
        return EXIT_SUCCESS;
    }
    bool ret=g_runner->stop();
    LOG("Delete last runner: id=%d\n",g_runner->getID());
    delete g_runner;
    g_runner= nullptr;
    if (ret) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}

const char * llama_gen(const char * prompt) {
    if (g_runner == nullptr) {
        LOG_ERR("Not init llama\n");
        return "";
    }
    std::string result = g_runner->generate(std::string(prompt));
    char* arr = new char[result.size() + 1];
    std::copy(result.begin(), result.end(), arr);
    arr[result.size()] = '\0';

    return arr;
}

const char * llama_chat(const char **roles,const char **contents, int size) {
    if (g_runner == nullptr) {
        LOG_ERR("Not init llama\n");
        return "";
    }
    std::vector<Message> msgs;

    for (int i = 0; i < size; i++) {
        Message msg;
        msg.role=roles[i];
        msg.content=contents[i];

        msgs.push_back(msg);
    }

    std::string result = g_runner->chat(msgs);
    char* arr = new char[result.size() + 1];
    std::copy(result.begin(), result.end(), arr);
    arr[result.size()] = '\0';

    return arr;
}

const char * whisper_gen(const char * model,const char * input) {
    WhisperService ws;

    std::string result = ws.generate(std::string(model),std::string(input));
    char* arr = new char[result.size() + 1];
    std::copy(result.begin(), result.end(), arr);
    arr[result.size()] = '\0';

    return arr;
}

int scheduler_start(const char * args) {
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

int scheduler_stop() {
    if (!Scheduler::instance().is_running()) {
        return EXIT_FAILURE;
    }
    if (Scheduler::instance().stop()) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}