#pragma once

#include <userver/clients/http/client.hpp>
#include <userver/clients/http/component.hpp>
#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>
#include <userver/components/raw_component_base.hpp>
#include <userver/yaml_config/merge_schemas.hpp>
#include <userver/yaml_config/schema.hpp>
#include <userver/yaml_config/yaml_config.hpp>

namespace Manager {
class HttpClient : public userver::components::ComponentBase {
public:
    static constexpr std::string_view kName = "client-http";

    HttpClient(const userver::components::ComponentConfig& config,
        const userver::components::ComponentContext& context);

    static userver::yaml_config::Schema GetStaticConfigSchema();

    std::shared_ptr<userver::clients::http::Response> CreateHttpRequest(std::string, const std::string&) const;

private:
    const std::string m_serviceUrl;
    userver::clients::http::Client& m_httpClient;
};

}

template <>
inline constexpr auto userver::components::kConfigFileMode<
    Manager::HttpClient> = userver::components::ConfigFileMode::kNotRequired;
