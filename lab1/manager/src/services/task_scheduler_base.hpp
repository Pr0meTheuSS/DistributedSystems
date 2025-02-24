#pragma once

#include <userver/components/loggable_component_base.hpp>

#include "components/http_worker_connection.hpp"
#include "models/worker.hpp"
#include "repositories/requests_repository.hpp"
#include "services/task_distribution_policy.hpp"

namespace Manager {

struct CrackRequest;

class TaskSchedulerBase {
public:
    TaskSchedulerBase(RequestsRepository&, TaskDistributor&, std::vector<Worker>);
    virtual ~TaskSchedulerBase();

    virtual CrackRequest Schedule(const CrackRequest&);
    virtual bool UpdateStatus(const CrackRequest&);
    virtual std::optional<CrackRequest> GetStatus(const std::string&) const;

private:
    RequestsRepository& m_repository;
    TaskDistributor& m_taskDistributionPolicy;
    std::vector<Worker> m_workers;
};

} // namespace Manager
