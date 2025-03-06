#pragma once

#include "models/sub_task.hpp"

namespace Manager {

class WorkerConnection {
public:
    ~WorkerConnection() = default;
    virtual bool Send(const SubTask&) = 0;
};

} // Manager
