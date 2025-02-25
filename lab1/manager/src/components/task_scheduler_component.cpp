#include "components/task_scheduler_component.hpp"

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include "clients/http_client.hpp"
#include "components/workers_factory_component.hpp"
#include "repositories/requests_repository.hpp"

namespace Manager {

TaskSchedulerComponent::TaskSchedulerComponent(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : TaskSchedulerBase(context.FindComponent<RequestsRepository>("repository-requests-in-memory"),
          context.FindComponent<WorkersFactoryComponent>("factory-workers").getWorkers())
    , userver::components::LoggableComponentBase(config, context)
{
}

} // namespace Manager
