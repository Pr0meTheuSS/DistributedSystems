#pragma once

#include <models/task.hpp>

namespace Manager {

class WorkerConnection {
public:
    ~WorkerConnection() = default;
    virtual bool Send(const Task&) = 0;
};

} // Manager
