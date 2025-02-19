#pragma once

#include <string>

namespace Manager {

enum class CrackStatus {
    IN_PROGRESS,
    READY,
    ERROR,
    UNDEFINED
};

struct CrackRequest final {
    std::string id;
    std::string hash;
    int maxLength;
    CrackStatus status;
};

} // namespace Manager

