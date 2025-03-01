#pragma once

#include <optional>
#include <string>

namespace Worker {

struct CrackRequest;
struct WorkerAnswer;

class RequestsRepository {
public:
    virtual std::optional<CrackRequest> GetByUUID(const std::string&) const = 0;
    virtual std::optional<CrackRequest> Create(const CrackRequest&) = 0;
    virtual std::optional<CrackRequest> UpdateStatus(const CrackRequest&) = 0;
    virtual std::optional<WorkerAnswer> GetAnswerByUUID(const std::string&) const = 0;
    virtual std::optional<WorkerAnswer> SaveAnswer(const WorkerAnswer&) = 0;

    virtual ~RequestsRepository() = default;
};

} // namespace Manager
