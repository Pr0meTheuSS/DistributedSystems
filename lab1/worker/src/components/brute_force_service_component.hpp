#pragma once

#include <chrono>
#include <iostream>
#include <memory>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>
#include <userver/utils/periodic_task.hpp>

#include <libenvpp/env.hpp>

#include "components/http_manager_connection.hpp"
#include "services/brute_force_service.hpp"

namespace Worker {

class BruteForceServiceComponent final : public BruteForceService,
                                         public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "component-brute-force-service";

    // // Component is valid after construction and is able to accept requests
    BruteForceServiceComponent(const userver::components::ComponentConfig& config,
        const userver::components::ComponentContext& context)
        : BruteForceService("abcdefghijklmnopqrstuvwxyz0123456789", std::cout)
        , userver::components::LoggableComponentBase(config, context)
        , m_httpManagerConnection(context.FindComponent<HttpManagerConnection>("component-http-manager-connection"))
    {
        // TODO: parse envs
        m_httpManagerConnection.SetUrl("http://localhost:8080");
    }

    virtual WorkerAnswer process(const SubTask& subTask) override
    {
        using namespace std::chrono_literals;
        userver::utils::PeriodicTask::Settings taskSettings(10s);

        auto callback = [this, s = std::move(subTask)]() mutable {
            m_httpManagerConnection.sendPing(s);
        };

        userver::utils::PeriodicTask pingTask("ping-manager", taskSettings, callback);

        WorkerAnswer answer = BruteForceService::process(subTask);
        m_httpManagerConnection.Send(answer);
        pingTask.Stop();
        return answer;
    }

    ~BruteForceServiceComponent() override = default;

private:
    HttpManagerConnection& m_httpManagerConnection;
};

} // namespace Worker
