#include "services/task_scheduler_base.hpp"

#include <algorithm>
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

TaskSchedulerBase::TaskSchedulerBase(
    RequestsRepository& repository, TasksRepository& tasksRepository, std::vector<Worker> workers)
    : m_repository(repository)
    , m_tasksRepo(tasksRepository)
    , m_workers(workers)
{
}

TaskSchedulerBase::~TaskSchedulerBase() { }

CrackRequest TaskSchedulerBase::Schedule(const CrackRequest& crackRequest)
{
    CrackRequest newCrackRequest { .id = crackRequest.id,
        .hash = crackRequest.hash,
        .maxLength = crackRequest.maxLength,
        .status = CrackStatus::IN_PROGRESS };

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

void TaskSchedulerBase::setLastPingFromWorkerBySubtask(const SubTask& subTask)
{
    std::string key = subTask.requestId + "_" + std::to_string(subTask.partNumber);
    m_subTaskPings.insert({ key, std::chrono::system_clock::now() });
}

std::optional<std::chrono::high_resolution_clock::time_point> TaskSchedulerBase::getLastPingFromWorkerBySubtask(
    const std::string& id, std::size_t partNumber) const
{
    std::string key = id + "_" + std::to_string(partNumber);
    return m_subTaskPings.contains(key)
        ? std::optional<std::chrono::high_resolution_clock::time_point> { m_subTaskPings.at(key) }
        : std::nullopt;
}

std::vector<SubTask> TaskSchedulerBase::separateToSubTasks(const Task& task)
{
    std::size_t partCount = m_workers.size();
    std::size_t partNumber = 0;

    std::vector<SubTask> subTasks;
    subTasks.reserve(partCount);

    for (; partNumber < partCount; partNumber++) {
        SubTask subTask { .requestId = task.requestId,
            .hash = task.hash,
            .maxLength = task.maxLength,
            .partCount = partCount,
            .partNumber = partNumber };
        subTasks.emplace_back(std::move(subTask));
    }
    return subTasks;
}

bool TaskSchedulerBase::UpdateStatus(const CrackRequest& crackRequst)
{
    return m_repository.UpdateStatus(crackRequst).has_value();
}

bool TaskSchedulerBase::UpdateStatus(const std::string& id, const CrackStatus& status)
{
    return m_repository.UpdateStatus(id, status).has_value();
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
        return false;
    }

    for (auto* worker : readyWorkers) {
        auto taskOpt = m_tasksRepo.getNextSubTask();
        if (!taskOpt)
            break;

        m_taskWorkerMap[taskOpt->requestId + "_" + std::to_string(taskOpt->partNumber)] = worker;
        worker->process(*taskOpt);
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
        m_repository.UpdateStatus(requestId, CrackStatus::READY);
    }

    distributeForReadyWorkers();
}

std::optional<WorkerAnswer> TaskSchedulerBase::saveWorkerAnswer(const WorkerAnswer& answer)
{
    return m_repository.SaveAnswer(answer);
}

std::optional<WorkerAnswer> TaskSchedulerBase::getAnswer(const std::string& id)
{
    auto task = m_repository.GetByUUID(id);
    if (!task) {
        return std::nullopt;
    }
    if (m_tasksRepo.areAllSubTasksCompleted(id)) {
        return m_repository.GetAnswerByUUID(id);
    }

    using namespace std::chrono_literals;

    // TODO: проверить, что готовы все неотвалившиеся ответы и только тогда отправить
    for (std::size_t part = 0; part < m_tasksRepo.getPartsCount(id); part++) {
        auto lastPing = getLastPingFromWorkerBySubtask(id, m_tasksRepo.getPartsCount(id));
        if (!lastPing || lastPing.value() < std::chrono::high_resolution_clock::now() - 10s) {
            UpdateStatus(id, CrackStatus::READY_PARTIAL_ANSWER);
        }
    }

    return m_repository.GetAnswerByUUID(id);
}

} // namespace Manager

