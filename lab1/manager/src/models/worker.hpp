#pragma once

#include <memory>

#include "components/http_worker_connection.hpp"
#include "models/sub_task.hpp"

namespace Manager {

class Worker {
public:
    explicit Worker(std::shared_ptr<WorkerConnection>);

    virtual ~Worker() = default;
    virtual void process(const SubTask& task);
    virtual bool isReady() const;
    virtual void setReady();

private:
    std::shared_ptr<WorkerConnection> m_connection;
    bool m_isReady;
};

} // namespace Manager
