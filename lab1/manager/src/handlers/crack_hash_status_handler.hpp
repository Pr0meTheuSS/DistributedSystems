#pragma once

#include <userver/components/component_list.hpp>
#include <userver/server/handlers/http_handler_json_base.hpp>

#include <userver/utest/using_namespace_userver.hpp>

#include "components/task_scheduler_component.hpp"

namespace Manager {

class CrackHashStatusHandler final : public userver::server::handlers::HttpHandlerJsonBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "handler-crack-hash-status";

    // // Component is valid after construction and is able to accept requests
    CrackHashStatusHandler(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    userver::formats::json::Value HandleRequestJsonThrow(
        const userver::server::http::HttpRequest&,
        const userver::formats::json::Value&,
        userver::server::request::RequestContext&) const override;

private:
    TaskSchedulerComponent& m_taskScheduler;
};

} // namespace Manager
