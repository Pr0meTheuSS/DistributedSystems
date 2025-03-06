#pragma once

#include <userver/components/component_list.hpp>
#include <userver/concurrent/background_task_storage.hpp>
#include <userver/server/handlers/http_handler_json_base.hpp>
#include <userver/utest/using_namespace_userver.hpp>

#include "components/brute_force_service_component.hpp"

namespace Worker {

class HttpManagerConnection;

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

    ~CrackHashHandler()
    {
        m_backgroundTaskStorage.CancelAndWait();
    }

private:
    HttpManagerConnection& m_connection;
    mutable userver::concurrent::BackgroundTaskStorage m_backgroundTaskStorage;
};

} // namespace Worker
