#pragma once

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>

#include <libenvpp/env.hpp>

#include <models/worker_connection.hpp>
#include <services/round_robin_distributor.hpp>
#include <services/task_distribution_policy.hpp>

namespace Manager {

class RoundRobinDistributorComponent final : public RoundRobinDistributor,
                                             public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "service-round-robin-distributor";

    // // Component is valid after construction and is able to accept requests
    RoundRobinDistributorComponent(const userver::components::ComponentConfig& config,
        const userver::components::ComponentContext& context)
        : userver::components::LoggableComponentBase(config, context)
    {
    }

    ~RoundRobinDistributorComponent() override = default;
    using RoundRobinDistributor::distribute;
};

} // Manager
