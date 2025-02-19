#pragma once

#include <optional>
#include <string>

namespace Manager {

struct CrackRequest;

class RequestsRepository {
public:
    virtual std::optional<CrackRequest> GetByUUID(const std::string&) const = 0;
    virtual std::optional<CrackRequest> Create(const CrackRequest&) = 0;
    virtual std::optional<CrackRequest> UpdateStatus(const CrackRequest&) = 0;
    virtual ~RequestsRepository() { };
};

} // namespace Manager
