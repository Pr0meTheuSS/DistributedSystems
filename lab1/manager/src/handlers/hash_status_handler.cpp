#include "handlers/hash_status_handler.hpp"

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include "components/task_scheduler_component.hpp"
#include "dto/crack_hash_status_response.hpp"
#include "models/crack_request.hpp"

namespace Manager {

HashStatusHandler::HashStatusHandler(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::server::handlers::HttpHandlerJsonBase(config, context)
    , m_taskScheduler(context.FindComponent<TaskSchedulerComponent>()) { };

userver::formats::json::Value HashStatusHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request,
    const userver::formats::json::Value&,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kGet: {
        const std::string& requestId = request.GetArg("requestId");
        std::optional<CrackRequest> req = m_taskScheduler.GetStatus(requestId);
        if (!req) {
            request.SetResponseStatus(userver::server::http::HttpStatus::kNotFound);
            throw userver::server::handlers::ClientError(userver::server::handlers::ExternalBody {
                fmt::format("Cannot find requst with uuid: {}", requestId) });
        }

        request.SetResponseStatus(userver::server::http::HttpStatus::kOk);
        CrackHashStatusResponse response {
            .requestId = req.value().id,
            .status = req.value().status
        };
        return userver::formats::json::ValueBuilder(response).ExtractValue();
    }
    default: {
        throw userver::server::handlers::ClientError(userver::server::handlers::ExternalBody {
            fmt::format("Unsupported method {}", request.GetMethod()) });
    }
    }
}

} // namespace Manager
