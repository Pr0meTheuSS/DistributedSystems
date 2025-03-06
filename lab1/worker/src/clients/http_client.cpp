#include "clients/http_client.hpp"

namespace Worker {

HttpClient::HttpClient(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : ComponentBase { config, context }
    , m_serviceUrl(config["service-url"].As<std::string>())
    , m_httpClient(context.FindComponent<userver::components::HttpClient>().GetHttpClient())
{
}

userver::yaml_config::Schema HttpClient::GetStaticConfigSchema()
{
    return userver::yaml_config::MergeSchemas<ComponentBase>(R"(
                    type: object
                    description: http client
                    additionalProperties: false
                    properties:
                        service-url:
                            type: string
                            description: other microservice listener url
                        )");
}

std::shared_ptr<userver::clients::http::Response> HttpClient::CreateHttpRequest(std::string action, const std::string& url) const
{
    return m_httpClient.CreateRequest().url(url).post().data(std::move(action)).perform();
}

}
