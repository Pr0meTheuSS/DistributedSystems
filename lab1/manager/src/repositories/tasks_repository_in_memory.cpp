#include "repositories/tasks_repository_in_memory.hpp"

#include <userver/components/raw_component_base.hpp>

namespace Manager {

// Конструктор
TasksRepositoryInMemory::TasksRepositoryInMemory(
    const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::components::LoggableComponentBase(config, context)
{
}

TasksRepositoryInMemory::~TasksRepositoryInMemory() = default;

void TasksRepositoryInMemory::registerSubTasks(const std::vector<SubTask>& subTasks)
{
    std::lock_guard lock(m_mutex);
    for (const auto& subTask : subTasks) {
        taskQueue.push(subTask);
        m_taskCompletionMap[subTask.requestId].push_back(false);
    }
}

std::optional<SubTask> TasksRepositoryInMemory::getNextSubTask()
{
    std::lock_guard lock(m_mutex);
    if (taskQueue.empty()) {
        return std::nullopt;
    }
    SubTask nextTask = taskQueue.front();
    taskQueue.pop();
    return nextTask;
}

void TasksRepositoryInMemory::completeSubTask(const std::string& requestId, std::size_t partNumber)
{
    std::lock_guard lock(m_mutex);
    if (m_taskCompletionMap.find(requestId) != m_taskCompletionMap.end()) {
        m_taskCompletionMap[requestId][partNumber] = true;
    }
}

bool TasksRepositoryInMemory::areAllSubTasksCompleted(const std::string& requestId) const
{
    std::lock_guard lock(m_mutex);
    auto it = m_taskCompletionMap.find(requestId);
    if (it == m_taskCompletionMap.end()) {
        return false;
    }
    const auto& completionMap = it->second;
    return std::all_of(completionMap.begin(), completionMap.end(), [](bool completed) { return completed; });
}

std::size_t TasksRepositoryInMemory::getPartsCount(const std::string& requestId) const
{
    std::lock_guard lock(m_mutex);
    auto it = m_taskCompletionMap.find(requestId);
    if (it == m_taskCompletionMap.end()) {
        return 0;
    }
    const auto& completionMap = it->second;

    return completionMap.size();
}

} // namespace Manager

template <>
inline constexpr auto userver::components::kConfigFileMode<
    Manager::TasksRepositoryInMemory> = userver::components::ConfigFileMode::kNotRequired;
