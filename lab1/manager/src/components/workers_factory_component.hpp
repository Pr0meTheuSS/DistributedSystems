#pragma once

#include <memory>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>

#include <libenvpp/env.hpp>

#include <models/worker_connection.hpp>
#include <services/round_robin_distributor.hpp>
#include <services/task_distribution_policy.hpp>

namespace Manager {

class WorkersFactoryComponent final : public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "factory-workers";

    // // Component is valid after construction and is able to accept requests
    WorkersFactoryComponent(const userver::components::ComponentConfig& config,
        const userver::components::ComponentContext& context)
        : userver::components::LoggableComponentBase(config, context)
    {
        const auto workersAmount = env::get<unsigned int>("WORKERS_AMOUNT").value();

        for (auto i = 1u; i <= workersAmount; i++) {
            std::string workerUrl = env::get<std::string>("WORKER_URL_" + std::to_string(i)).value();

            auto connection = std::make_unique<HttpWorkerConnection>(config, context);
            connection->SetUrl(workerUrl);

            m_workers.push_back(Worker(std::move(connection))); // Передаём владение в Worker
        }
    }

    std::vector<Worker> getWorkers() const
    {
        return m_workers;
    }

    ~WorkersFactoryComponent() override = default;

private:
    std::vector<Worker> m_workers;
};

} // Manager

template <>
inline constexpr auto userver::components::kConfigFileMode<
    Manager::WorkersFactoryComponent> = userver::components::ConfigFileMode::kNotRequired;
