#pragma once

#include <unordered_map>

#include <userver/components/loggable_component_base.hpp>
#include <userver/engine/mutex.hpp>

#include "dto/worker_answer.hpp"
#include "models/crack_request.hpp"
#include "repositories/requests_repository.hpp"

namespace Manager {

class RequestsRepositoryInMemory final : public RequestsRepository,
                                         public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "repository-requests-in-memory";

    RequestsRepositoryInMemory(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    ~RequestsRepositoryInMemory() override;

    std::optional<CrackRequest> GetByUUID(const std::string& uuid) const override;
    std::optional<CrackRequest> Create(const CrackRequest& request) override;
    std::optional<CrackRequest> UpdateStatus(const CrackRequest&) override;

    std::optional<WorkerAnswer> GetAnswerByUUID(const std::string&) const override;
    std::optional<WorkerAnswer> SaveAnswer(const WorkerAnswer&) override;

private:
    std::unordered_map<std::string, CrackRequest> m_crackRequests;
    std::unordered_map<std::string, WorkerAnswer> m_workerAnswers;
    mutable userver::engine::Mutex m_mutex;
};

} // namespace Manager
