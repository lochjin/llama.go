#pragma once

#include <iostream>
class WhisperService {
private:

public:
    WhisperService();
    ~WhisperService();

    const std::string generate(const std::string& model,const std::string& input);
};