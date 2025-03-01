#pragma once

#include <userver/components/component_list.hpp>
#include <userver/server/handlers/http_handler_json_base.hpp>

#include <userver/utest/using_namespace_userver.hpp>

namespace Worker {

class CrackHashHandler final : public userver::server::handlers::HttpHandlerJsonBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "handler-hash-crack";

    // // Component is valid after construction and is able to accept requests
    CrackHashHandler(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    userver::formats::json::Value HandleRequestJsonThrow(
        const userver::server::http::HttpRequest&,
        const userver::formats::json::Value&,
        userver::server::request::RequestContext&) const override;

    // private:
    // WorkerServiceComponent& m_workerService;
};

} // namespace Worker
