#include <iostream>
#include <thread>
#include <chrono>
#include <cstdlib>
#include <future>
#include <vector>

int main() {
    std::string model("/Users/jin/temp/ggufs/DeepSeek-8B-q8_0.gguf");

    std::vector<std::string> v_args;
    v_args.push_back("test_runner_gen");
    v_args.push_back("-m");
    v_args.push_back(model);
    v_args.push_back("--seed");
    v_args.push_back("0");



    std::cout<<"stop:"<<std::endl;
    return EXIT_SUCCESS;
}