#pragma once

#include <string>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>

#include "clients/http_client.hpp"
#include "dto/worker_answer.hpp"
#include "models/sub_task.hpp"

namespace Worker {

class HttpManagerConnection final : public userver::components::LoggableComponentBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "component-http-manager-connection";

    // Component is valid after construction and is able to accept requests
    HttpManagerConnection(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    ~HttpManagerConnection() override = default;
    bool Send(const WorkerAnswer&);

    void sendPing(const SubTask&);

    void SetUrl(std::string);

private:
    HttpClient& m_httpClient;
    std::string m_url;
};

} // Manager
