#include "handlers/crack_hash_handler.hpp"

#include <iostream>

#include <fmt/format.h>
#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>
#include <userver/utils/async.hpp>

#include "models/sub_task.hpp"
#include "services/bf_service.hpp"

namespace Worker {

CrackHashHandler::CrackHashHandler(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::server::handlers::HttpHandlerJsonBase(config, context) { };

userver::formats::json::Value CrackHashHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request,
    const userver::formats::json::Value& requestBody,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kPost: {
        auto subTask = requestBody.As<SubTask>();

        WorkerService workerService("abcdefghijklmnopqrstuvwxyz0123456789");
        // workerService.process(subTask);

        // TODO: run coroutine without wait
        auto task = userver::utils::Async("brutforce", [&]() {
            workerService.process(subTask);
        });
        task.Wait();

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
