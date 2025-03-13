#include "handlers/crack_hash_handler.hpp"

#include <chrono>
#include <iostream>

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>
#include <userver/utils/async.hpp>

#include <fmt/format.h>

#include "components/http_manager_connection.hpp"
#include "models/sub_task.hpp"
#include "services/brute_force_service.hpp"

namespace Worker {

CrackHashHandler::CrackHashHandler(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::server::handlers::HttpHandlerJsonBase(config, context)
    , m_connection(context.FindComponent<HttpManagerConnection>("component-http-manager-connection"))
    , m_backgroundTaskStorage(context.GetTaskProcessor("workers-task-processor")) {
    };

userver::formats::json::Value CrackHashHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request,
    const userver::formats::json::Value& requestBody,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kPost: {
        auto subTask = requestBody.As<SubTask>();
        BruteForceService workerService("abcdefghijklmnopqrstuvwxyz0123456789", std::cout);

        m_backgroundTaskStorage.AsyncDetach("bf", [this, w = std::move(workerService), s = std::move(subTask)]() mutable {
            using namespace std::chrono_literals;
            userver::utils::PeriodicTask::Settings taskSettings(1s, { userver::utils::PeriodicTask::Flags::kNow });

            auto callback = [&]() mutable {
                m_connection.sendPing(s);
            };

            userver::utils::PeriodicTask pingTask("ping-manager", taskSettings, callback);

            WorkerAnswer answer = w.process(s);
            m_connection.Send(answer);
            pingTask.Stop();
        });

        request.SetResponseStatus(userver::server::http::HttpStatus::kNoContent);
        return {};
    }
    default: {
        throw userver::server::handlers::ClientError(userver::server::handlers::ExternalBody {
            fmt::format("Unsupported method {}", request.GetMethod()) });
    }
    }
}

} // namespace Manager
