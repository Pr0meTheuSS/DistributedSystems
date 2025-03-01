#pragma once

#include <string>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>

#include "clients/http_client.hpp"
#include "models/worker_connection.hpp"

namespace Worker {

class HttpWorkerConnection final : public WorkerConnection,
                                   public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "service-http-worker-connection";

    // Component is valid after construction and is able to accept requests
    HttpWorkerConnection(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    ~HttpWorkerConnection() override = default;
    bool Send(const SubTask&) override;

    void SetUrl(std::string url)
    {
        m_url = url;
    }

private:
    HttpClient& m_httpClient;
    std::string m_url;
};

} // Manager
