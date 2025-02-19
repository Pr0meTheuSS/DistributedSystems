#include <repositories/requests_repository_in_memory.hpp>

#include <iostream>

#include <fmt/format.h>
#include <sole.fwd.hpp>
#include <sole.hpp>

#include <models/crack_request.hpp>

namespace Manager {

RequestsRepositoryInMemory::RequestsRepositoryInMemory(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::components::LoggableComponentBase(config, context) { };

RequestsRepositoryInMemory::~RequestsRepositoryInMemory() { };

std::optional<CrackRequest> RequestsRepositoryInMemory::GetByUUID(const std::string& uuid) const
{
    return m_crackRequests.contains(uuid) ? std::optional { m_crackRequests.at(uuid) } : std::nullopt;
}

std::optional<CrackRequest> RequestsRepositoryInMemory::Create(const CrackRequest& request)
{
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
    if (m_crackRequests.contains(crackRequest.id)) {
        m_crackRequests[crackRequest.id].status = crackRequest.status;
        return m_crackRequests[crackRequest.id];
    }
    return std::nullopt;
}

} // namespace Managerauto
