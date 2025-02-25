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
    std::size_t maxLength;
    CrackStatus status;
};

} // namespace Manager

