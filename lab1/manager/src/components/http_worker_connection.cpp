#include "components/http_worker_connection.hpp"

#include <exception>
#include <iostream>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>
#include <userver/formats/json/value_builder.hpp>

#include "models/sub_task.hpp"

namespace Manager {

HttpWorkerConnection::HttpWorkerConnection(
    const userver::components::ComponentConfig& config, const userver::components::ComponentContext& context)
    : LoggableComponentBase(config, context)
    , m_httpClient(context.FindComponent<HttpClient>())
{
}

bool HttpWorkerConnection::Send(const SubTask& task)
{
    auto taskJson = userver::formats::json::ValueBuilder(task).ExtractValue();
    const auto requestPath = this->m_url + "/api/hash/crack";
    const auto requestBody = userver::formats::json::ToStableString(taskJson);
    try {
        const auto response = m_httpClient.CreateHttpRequest(requestBody, requestPath);
    } catch (const std::exception& e) {
        std::cerr << fmt::format("Cannot send http request to {}. \n Request body: {}\n", requestPath, requestBody);
    }
    return true;
}

} // namespace Manager