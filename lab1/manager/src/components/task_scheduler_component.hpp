#pragma once

#include <userver/components/loggable_component_base.hpp>

#include "clients/http_client.hpp"
#include "services/task_scheduler_base.hpp"

namespace Manager {

struct CrackRequest;

class TaskSchedulerComponent final : public TaskSchedulerBase,
                                     public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "service-task-scheduler";

    TaskSchedulerComponent(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    using TaskSchedulerBase::GetStatus;
    using TaskSchedulerBase::Schedule;
    using TaskSchedulerBase::UpdateStatus;
};

} // namespace Manager
