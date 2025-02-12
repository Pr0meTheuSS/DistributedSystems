#pragma once

#include <userver/components/component_list.hpp>
#include <userver/server/handlers/http_handler_json_base.hpp>

// Note: this is for the purposes of tests/samples only
#include <userver/utest/using_namespace_userver.hpp>

namespace samples::hello {

class HashCrackHandler final : public userver::server::handlers::HttpHandlerJsonBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "handler-hash-crack";

    // // Component is valid after construction and is able to accept requests
    // using HttpHandlerBase::HttpHandlerBase;
    HashCrackHandler(const userver::components::ComponentConfig& config,
        const userver::components::ComponentContext& context)
        : userver::server::handlers::HttpHandlerJsonBase(config, context) { };

    userver::formats::json::Value HandleRequestJsonThrow(
        const userver::server::http::HttpRequest&,
        const userver::formats::json::Value&,
        userver::server::request::RequestContext&) const override;
};

} // namespace samples::hello