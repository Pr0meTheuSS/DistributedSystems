#include "components/http_manager_connection.hpp"

#include <exception>

#include <userver/clients/http/component.hpp>
#include <userver/components/loggable_component_base.hpp>
#include <userver/formats/json/value_builder.hpp>

#include "models/sub_task.hpp"

namespace Worker {

HttpManagerConnection::HttpManagerConnection(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : LoggableComponentBase(config, context)
    , m_httpClient(context.FindComponent<HttpClient>())
{
}

bool HttpManagerConnection::Send(const WorkerAnswer& answer)
{
    auto answerJson = userver::formats::json::ValueBuilder(answer).ExtractValue();
    try {
        const auto response = m_httpClient.CreateHttpRequest(
            userver::formats::json::ToStableString(answerJson),
            m_url + "/internal/api/manager/hash/crack/request");
    } catch (const std::exception&) {
        return false;
    }
    return true;
}

void HttpManagerConnection::sendPing(const SubTask& subTask)
{
    auto subTaskJson = userver::formats::json::ValueBuilder(subTask).ExtractValue();
    try {
        const auto response = m_httpClient.CreateHttpRequest(
            userver::formats::json::ToStableString(subTaskJson),
            this->m_url + "/internal/api/manager/subtask/status");
    } catch (const std::exception&) {
    }
}

void HttpManagerConnection::SetUrl(std::string url)
{
    m_url = url;
}

} // namespace Manager
