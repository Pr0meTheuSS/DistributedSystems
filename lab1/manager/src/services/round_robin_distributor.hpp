#pragma once
#include "services/task_distribution_policy.hpp"

#include "models/task.hpp"
#include "models/worker.hpp"

namespace Manager {

class RoundRobinDistributor : public TaskDistributor {
public:
    RoundRobinDistributor()
    {
    }

    virtual ~RoundRobinDistributor() = default;

    virtual void distribute(const Task& task, const std::vector<Worker>& workers) override
    {
        if (workers.empty()) {
            return;
        }

        std::size_t offset = 0;
        for (auto worker : workers) {
            Task subTask {
                .hash = task.hash,
                .requestId = task.requestId,
                .maxLength = task.maxLength,
                .offset = offset++,
                .size = workers.size()
            };

            worker.Process(subTask);
        }
    }
};

}
