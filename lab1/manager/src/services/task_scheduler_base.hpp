#pragma once

#include <queue>
#include <unordered_map>
#include <vector>

#include <userver/components/loggable_component_base.hpp>
#include <userver/engine/mutex.hpp>

#include "components/http_worker_connection.hpp"
#include "dto/worker_answer.hpp"
#include "models/worker.hpp"
#include "repositories/requests_repository.hpp"
#include "services/task_distribution_policy.hpp"

namespace Manager {

struct CrackRequest;
struct SubTask;

class TaskSchedulerBase {
public:
    TaskSchedulerBase(RequestsRepository&, std::vector<Worker>);
    virtual ~TaskSchedulerBase();

    virtual CrackRequest Schedule(const CrackRequest&);
    virtual bool UpdateStatus(const CrackRequest&);
    virtual void completeSubTask(const std::string&, std::size_t);
    virtual std::optional<CrackRequest> GetStatus(const std::string&) const;

    virtual std::optional<WorkerAnswer> saveWorkerAnswer(const WorkerAnswer&);
    virtual std::optional<WorkerAnswer> getAnswer(const std::string&) const;

private:
    std::vector<SubTask> separateToSubTasks(const Task& task);
    void registerSubTasks(const std::vector<SubTask>& subTasks);
    void distributeForReadyWorkers();

    RequestsRepository& m_repository;
    std::vector<Worker> m_workers;

    // TODO: abstraction leak
    userver::engine::Mutex m_mutex;

    std::unordered_map<std::string, std::vector<bool>> taskCompletionMap;
    std::unordered_map<std::string, Worker*> taskWorkerMap;
    std::queue<SubTask> taskQueue;
};

} // namespace Manager
