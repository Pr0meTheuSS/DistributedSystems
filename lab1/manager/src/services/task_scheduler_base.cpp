#include "services/task_scheduler_base.hpp"

#include <sole.fwd.hpp>
#include <sole.hpp>

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include <dto/crack_hash_request.hpp>
#include <dto/crack_hash_response.hpp>
#include <models/crack_request.hpp>
#include <models/task.hpp>
#include <repositories/requests_repository.hpp>
#include <repositories/requests_repository_in_memory.hpp>
#include <services/task_distribution_policy.hpp>

namespace Manager {

TaskSchedulerBase::TaskSchedulerBase(RequestsRepository& repository,
    TaskDistributor& taskDistributionPolicy,
    std::vector<Worker> workers)
    : m_repository(repository)
    , m_taskDistributionPolicy(taskDistributionPolicy)
    , m_workers(workers)
{
}

TaskSchedulerBase::~TaskSchedulerBase() { }

CrackRequest TaskSchedulerBase::Schedule(const CrackRequest& crackRequest)
{
    CrackRequest newCrackRequest {
        .id = crackRequest.id,
        .hash = crackRequest.hash,
        .maxLength = crackRequest.maxLength,
        .status = CrackStatus::IN_PROGRESS
    };

    auto created = m_repository.Create(newCrackRequest);
    if (!created) {
        return {};
    }

    Task task {
        .hash = created.value().hash,
        .requestId = created.value().id,
        .maxLength = created.value().maxLength,
    };

    m_taskDistributionPolicy.distribute(task, m_workers);
    return created.value();
}

bool TaskSchedulerBase::UpdateStatus(const CrackRequest& crackRequst)
{
    return m_repository.UpdateStatus(crackRequst).has_value();
}

std::optional<CrackRequest> TaskSchedulerBase::GetStatus(const std::string& id) const
{
    return m_repository.GetByUUID(id);
}

} // namespace Manager
