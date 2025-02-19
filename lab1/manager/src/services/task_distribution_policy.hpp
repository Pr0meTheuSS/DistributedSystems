#pragma once

#include <vector>

#include <models/task.hpp>
#include <models/worker.hpp>

namespace Manager {

class TaskDistributor {
public:
    virtual ~TaskDistributor() = default;
    // TODO: make able to call with all iterable containers
    virtual void distribute(const Task& task, const std::vector<Worker>& workers) = 0;
};

} // Manager