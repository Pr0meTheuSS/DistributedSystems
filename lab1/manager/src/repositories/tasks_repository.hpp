#pragma once

#include <optional>
#include <string>
#include <vector>

#include "models/sub_task.hpp"

namespace Manager {

// Интерфейс репозитория для задач
class TasksRepository {
public:
    virtual ~TasksRepository() = default;

    virtual void registerSubTasks(const std::vector<SubTask>& subTasks) = 0;
    virtual std::optional<SubTask> getNextSubTask() = 0;
    virtual void completeSubTask(const std::string& requestId, std::size_t partNumber) = 0;
    virtual bool areAllSubTasksCompleted(const std::string& requestId) const = 0;
};

} // namespace Manager
