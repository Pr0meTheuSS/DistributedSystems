#pragma once

#include <memory>
#include <string>
#include <vector>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>

#include <libenvpp/env.hpp>

#include "models/worker.hpp"
#include "models/worker_connection.hpp"

namespace Manager {

class WorkersFactoryComponent final : public userver::components::LoggableComponentBase {
public:
    static constexpr std::string_view kName = "factory-workers";

    WorkersFactoryComponent(
        const userver::components::ComponentConfig& config, const userver::components::ComponentContext& context);

    std::vector<Worker> getWorkers() const;

    ~WorkersFactoryComponent() override = default;

private:
    std::vector<Worker> m_workers;
};

} // namespace Manager

template <> inline constexpr auto userver::components::kConfigFileMode<Manager::WorkersFactoryComponent> =
    userver::components::ConfigFileMode::kNotRequired;
