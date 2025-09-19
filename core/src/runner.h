#pragma once

#include "event_processor.h"
#include "sampling.h"
#include "message.h"
#include "singleton.h"

class Runner: public patterns::Singleton<Runner> {
    friend class patterns::Singleton<Runner>;

private:
    int m_id;
    std::vector<std::string> m_args;
    EventProcessor m_eprocessor;
    std::atomic<bool> m_running;
    bool m_async;

    llama_context           * m_ctx;
    llama_model             * m_model;
    common_sampler          * m_smpl;
    common_params           * m_params;
    std::string               m_prompt;

    std::vector<llama_token> * m_input_tokens;
    std::ostringstream       * m_output_ss;
    std::vector<llama_token> * m_output_tokens;

    Runner();
    ~Runner();
public:
    bool start(int id,const std::vector<std::string>& args,bool async= false,const std::string& prt="");
    bool stop();
    const std::string generate(const std::string& prompt);
    const std::string chat(const std::vector<Message>& mgs);
    int getID();
    bool isRunning();

    bool getPrompt(EventProcessor::Event& event);
};