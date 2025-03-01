#include "services/task_scheduler_base.hpp"

#include <algorithm>
#include <iostream>
#include <optional>
#include <ranges>
#include <vector>

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include "dto/crack_hash_request.hpp"
#include "dto/crack_hash_response.hpp"
#include "models/crack_request.hpp"
#include "models/task.hpp"
#include "repositories/requests_repository.hpp"
#include "repositories/requests_repository_in_memory.hpp"

namespace Manager {

TaskSchedulerBase::TaskSchedulerBase(RequestsRepository& repository,
    TasksRepository& tasksRepository,
    std::vector<Worker> workers)
    : m_repository(repository)
    , m_tasksRepo(tasksRepository)
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
        newCrackRequest.status = CrackStatus::ERROR;
        return newCrackRequest;
    }
    newCrackRequest = created.value();

    Task task {
        .requestId = created.value().id,
        .hash = created.value().hash,
        .maxLength = created.value().maxLength,
    };

    std::vector<SubTask> subTasks = separateToSubTasks(task);
    m_tasksRepo.registerSubTasks(subTasks);
    if (!distributeForReadyWorkers()) {
        newCrackRequest.status = CrackStatus::IN_QUEUE;
        return newCrackRequest;
    }

    return newCrackRequest;
}

std::vector<SubTask> TaskSchedulerBase::separateToSubTasks(const Task& task)
{
    std::size_t partCount = m_workers.size();
    std::size_t partNumber = 0;

    std::vector<SubTask> subTasks;
    subTasks.reserve(partCount);

    for (; partNumber < partCount; partNumber++) {
        SubTask subTask {
            .requestId = task.requestId,
            .hash = task.hash,
            .maxLength = task.maxLength,
            .partCount = partCount,
            .partNumber = partNumber
        };
        subTasks.emplace_back(std::move(subTask));
    }
    return subTasks;
}

bool TaskSchedulerBase::UpdateStatus(const CrackRequest& crackRequst)
{
    return m_repository.UpdateStatus(crackRequst).has_value();
}

std::optional<CrackRequest> TaskSchedulerBase::GetStatus(const std::string& id) const
{
    return m_repository.GetByUUID(id);
}

std::vector<Worker*> TaskSchedulerBase::getReadyWorkers()
{
    std::vector<Worker*> readyWorkers;
    for (auto& worker : m_workers) {
        if (worker.isReady()) {
            readyWorkers.push_back(&worker);
        }
    }
    return readyWorkers;
}

bool TaskSchedulerBase::distributeForReadyWorkers()
{
    auto readyWorkers = getReadyWorkers();
    if (readyWorkers.empty()) {
        std::cout << "[INFO] Нет свободных воркеров\n";
        return false;
    }

    for (auto* worker : readyWorkers) {
        auto taskOpt = m_tasksRepo.getNextSubTask();
        if (!taskOpt)
            break;

        m_taskWorkerMap[taskOpt->requestId + "_" + std::to_string(taskOpt->partNumber)] = worker;
        worker->Process(*taskOpt);
    }

    return true;
}

void TaskSchedulerBase::completeSubTask(const std::string& requestId, std::size_t partNumber)
{
    m_tasksRepo.completeSubTask(requestId, partNumber);

    std::string taskKey = requestId + "_" + std::to_string(partNumber);
    auto workerIt = m_taskWorkerMap.find(taskKey);
    if (workerIt != m_taskWorkerMap.end()) {
        workerIt->second->setReady();
        m_taskWorkerMap.erase(workerIt);
    }

    if (!m_tasksRepo.areAllSubTasksCompleted(requestId)) {
        distributeForReadyWorkers();
    }
}

std::optional<WorkerAnswer> TaskSchedulerBase::saveWorkerAnswer(const WorkerAnswer& answer)
{
    return m_repository.SaveAnswer(answer);
}

std::optional<WorkerAnswer> TaskSchedulerBase::getAnswer(const std::string& id) const
{
    return m_repository.GetAnswerByUUID(id);
}

} // namespace Manager

