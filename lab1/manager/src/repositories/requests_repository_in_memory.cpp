#include "repositories/requests_repository_in_memory.hpp"

#include <iostream>

#include <fmt/format.h>
#include <sole.fwd.hpp>
#include <sole.hpp>

#include "models/crack_request.hpp"

namespace Manager {

RequestsRepositoryInMemory::RequestsRepositoryInMemory(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::components::LoggableComponentBase(config, context) { };

RequestsRepositoryInMemory::~RequestsRepositoryInMemory() = default;

std::optional<CrackRequest> RequestsRepositoryInMemory::GetByUUID(const std::string& uuid) const
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);

    return m_crackRequests.contains(uuid) ? std::optional { m_crackRequests.at(uuid) } : std::nullopt;
}

std::optional<CrackRequest> RequestsRepositoryInMemory::Create(const CrackRequest& request)
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);

    CrackRequest newCrackRequest {
        .id = sole::uuid0().str(),
        .hash = request.hash,
        .maxLength = request.maxLength,
        .status = request.status
    };

    auto [inserted, ok] = m_crackRequests.insert({ newCrackRequest.id, newCrackRequest });
    return ok ? std::optional<CrackRequest>(inserted->second) : std::nullopt;
}

std::optional<CrackRequest> RequestsRepositoryInMemory::UpdateStatus(const CrackRequest& crackRequest)
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);

    if (m_crackRequests.contains(crackRequest.id)) {
        m_crackRequests[crackRequest.id].status = crackRequest.status;
        return m_crackRequests[crackRequest.id];
    }
    return std::nullopt;
}

std::optional<WorkerAnswer> RequestsRepositoryInMemory::GetAnswerByUUID(const std::string& id) const
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);
    return m_workerAnswers.contains(id) ? std::optional { m_workerAnswers.at(id) } : std::nullopt;
}

std::optional<WorkerAnswer> RequestsRepositoryInMemory::SaveAnswer(const WorkerAnswer& answer)
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);
    auto [inserted, ok] = m_workerAnswers.insert({ answer.requestId, answer });
    return ok ? std::optional<WorkerAnswer>(inserted->second) : std::nullopt;
}

} // namespace Managerauto
