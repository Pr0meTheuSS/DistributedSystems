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
#include "services/task_distribution_policy.hpp"

namespace Manager {

TaskSchedulerBase::TaskSchedulerBase(RequestsRepository& repository,
    std::vector<Worker> workers)
    : m_repository(repository)
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

    Task task {
        .requestId = created.value().id,
        .hash = created.value().hash,
        .maxLength = created.value().maxLength,
    };

    std::vector<SubTask> subTasks = separateToSubTasks(task);
    registerSubTasks(subTasks);
    distributeForReadyWorkers();

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

std::vector<SubTask> TaskSchedulerBase::separateToSubTasks(const Task& task)
{
    size_t partCount = m_workers.size();
    std::vector<SubTask> subTasks;
    subTasks.reserve(partCount);

    std::cout << "[INFO] Разделяем задачу " << task.requestId << " на " << partCount << " частей\n";

    for (size_t i = 0; i < partCount; ++i) {
        subTasks.push_back({
            .requestId = task.requestId,
            .hash = task.hash,
            .maxLength = task.maxLength,
            .partCount = partCount,
            .partNumber = i,
        });
        std::cout << "[INFO] Создана подзадача " << task.requestId << " #" << i << "\n";
    }
    return subTasks;
}

// Регистрация тасок в системе
void TaskSchedulerBase::registerSubTasks(const std::vector<SubTask>& subTasks)
{
    if (subTasks.empty())
        return;

    std::lock_guard<userver::engine::Mutex> lock(m_mutex);
    taskCompletionMap[subTasks[0].requestId] = std::vector<bool>(subTasks.size(), false);

    for (const auto& subTask : subTasks) {
        taskQueue.push(subTask);
    }

    std::cout << "[INFO] Зарегистрированы подзадачи для " << subTasks[0].requestId
              << " (всего " << subTasks.size() << " частей), добавлены в очередь\n";
}

// Распределение задач по свободным воркерам
void TaskSchedulerBase::distributeForReadyWorkers()
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);

    std::cout << "[INFO] Начинаем распределение задач по свободным воркерам\n";

    for (auto& worker : m_workers) {
        if (!worker.isReady() || taskQueue.empty())
            continue;

        SubTask task = taskQueue.front();
        taskQueue.pop();

        worker.Process(task);
        taskWorkerMap[task.requestId + "_" + std::to_string(task.partNumber)] = &worker;

        std::cout << "[INFO] Воркер " << &worker << " взял подзадачу " << task.requestId
                  << " #" << task.partNumber << " в работу\n";
    }
}

// Завершение подзадачи
void TaskSchedulerBase::completeSubTask(const std::string& requestId, std::size_t partNumber)
{
    std::lock_guard<userver::engine::Mutex> lock(m_mutex);

    auto it = taskCompletionMap.find(requestId);
    if (it == taskCompletionMap.end())
        return;

    it->second[partNumber] = true;
    std::cout << "[INFO] Завершена подзадача " << requestId << " #" << partNumber << "\n";

    if (std::all_of(it->second.begin(), it->second.end(), [](bool completed) { return completed; })) {
        std::cout << "[INFO] Все подзадачи " << requestId << " завершены. Очищаем мета-данные\n";

        taskCompletionMap.erase(it);
        taskWorkerMap.erase(requestId);
    }

    std::string taskKey = requestId + "_" + std::to_string(partNumber);
    auto workerIt = taskWorkerMap.find(taskKey);
    if (workerIt != taskWorkerMap.end()) {
        workerIt->second->setReady();
        taskWorkerMap.erase(workerIt);
        std::cout << "[INFO] Воркер освободился, ищем новую задачу\n";
        distributeForReadyWorkers();
    }
}

} // namespace Manager

