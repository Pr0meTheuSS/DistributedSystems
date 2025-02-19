#pragma once

#include <userver/components/component_list.hpp>
#include <userver/server/handlers/http_handler_json_base.hpp>

#include <components/task_scheduler_component.hpp>

namespace Manager {

class HashStatusHandler final : public userver::server::handlers::HttpHandlerJsonBase {
public:
    // `kName` is used as the component name in static config
    static constexpr std::string_view kName = "handler-hash-status";

    // // Component is valid after construction and is able to accept requests
    HashStatusHandler(const userver::components::ComponentConfig&,
        const userver::components::ComponentContext&);

    userver::formats::json::Value HandleRequestJsonThrow(
        const userver::server::http::HttpRequest&,
        const userver::formats::json::Value&,
        userver::server::request::RequestContext&) const override;

private:
    TaskSchedulerComponent& m_taskScheduler;
};

} // namespace Manager
