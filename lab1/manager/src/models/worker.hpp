#pragma once

#include <memory>

#include "components/http_worker_connection.hpp"
#include "models/task.hpp"

namespace Manager {

class Worker {
public:
    explicit Worker(std::shared_ptr<WorkerConnection> connection);

    virtual ~Worker() = default;
    virtual void Process(const Task& task);

private:
    std::shared_ptr<WorkerConnection> m_connection;
};

} // Manager
