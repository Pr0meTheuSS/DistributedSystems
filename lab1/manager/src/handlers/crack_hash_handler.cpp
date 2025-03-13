#include "handlers/crack_hash_handler.hpp"

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include <fmt/format.h>

#include "dto/crack_hash_request.hpp"
#include "dto/crack_hash_response.hpp"
#include "models/crack_request.hpp"

namespace Manager {

CrackHashHandler::CrackHashHandler(
    const userver::components::ComponentConfig& config, const userver::components::ComponentContext& context)
    : userver::server::handlers::HttpHandlerJsonBase(config, context)
    , m_taskScheduler(context.FindComponent<TaskSchedulerComponent>()) { };

userver::formats::json::Value CrackHashHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request, const userver::formats::json::Value& requestBody,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kPost: {
        auto req = requestBody.As<CrackHashRequest>();
        CrackRequest crackRequest {
            .id = "", .hash = req.hash, .maxLength = req.maxLength, .status = CrackStatus::UNDEFINED
        };

        crackRequest = m_taskScheduler.Schedule(crackRequest);
        CrackHashResponse response { crackRequest.id };
        request.SetResponseStatus(userver::server::http::HttpStatus::kCreated);
        return userver::formats::json::ValueBuilder(response).ExtractValue();
    }
    default: {
        throw userver::server::handlers::ClientError(
            userver::server::handlers::ExternalBody { fmt::format("Unsupported method {}", request.GetMethod()) });
    }
    }
}

} // namespace Manager
