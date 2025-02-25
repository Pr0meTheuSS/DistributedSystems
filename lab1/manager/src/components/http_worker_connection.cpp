#include "components/http_worker_connection.hpp"

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>
#include <userver/formats/json/value_builder.hpp>

#include "models/sub_task.hpp"

namespace Manager {

HttpWorkerConnection::HttpWorkerConnection(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : LoggableComponentBase(config, context)
    , m_httpClient(context.FindComponent<HttpClient>())
{
}

bool HttpWorkerConnection::Send(const SubTask& task)
{
    auto taskJson = userver::formats::json::ValueBuilder(task).ExtractValue();

    const auto response = m_httpClient.CreateHttpRequest(
        userver::formats::json::ToStableString(taskJson),
        this->m_url);
    return true;
}

} // namespace Manager