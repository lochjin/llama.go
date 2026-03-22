//
// Created by jinjin on 2026/3/21.
//

#pragma once
#include <vector>
#include <string>
#include "server-context.h"

class Server {
private:
    std::unique_ptr<server_routes> routes;
    server_context ctx_server;

public:
    Server();
    ~Server();

    bool start(const std::vector<std::string>& args);
    bool stop();

    void process();
};