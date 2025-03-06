#pragma once

#include <chrono>
#include <queue>
#include <unordered_map>

#include <userver/components/loggable_component_base.hpp>
#include <userver/engine/mutex.hpp>

#include "models/sub_task.hpp"
#include "repositories/tasks_repository.hpp"

namespace Manager {

class TasksRepositoryInMemory final : public TasksRepository, public userver::components::LoggableComponentBase {
public:
    static constexpr std::string_view kName = "repository-tasks-in-memory";

    TasksRepositoryInMemory(const userver::components::ComponentConfig&, const userver::components::ComponentContext&);
    ~TasksRepositoryInMemory() override;

    void registerSubTasks(const std::vector<SubTask>& subTasks) override;
    std::optional<SubTask> getNextSubTask() override;
    void completeSubTask(const std::string& requestId, std::size_t partNumber) override;
    bool areAllSubTasksCompleted(const std::string& requestId) const override;

    std::size_t getPartsCount(const std::string&) const override;

private:
    std::unordered_map<std::string, std::vector<bool>> m_taskCompletionMap;
    std::queue<SubTask> taskQueue;
    mutable userver::engine::Mutex m_mutex;
};

} // namespace Manager

