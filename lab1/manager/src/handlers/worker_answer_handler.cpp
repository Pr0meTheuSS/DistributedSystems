#include "handlers/worker_answer_handler.hpp"

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>
#include <userver/formats/serialize/common_containers.hpp>

#include <fmt/format.h>

#include "dto/crack_hash_request.hpp"
#include "dto/crack_hash_response.hpp"
#include "dto/worker_answer.hpp"
#include "models/crack_request.hpp"

namespace Manager {

WorkerAnswerHandler::WorkerAnswerHandler(
    const userver::components::ComponentConfig& config, const userver::components::ComponentContext& context)
    : userver::server::handlers::HttpHandlerJsonBase(config, context)
    , m_taskScheduler(context.FindComponent<TaskSchedulerComponent>()) { };

userver::formats::json::Value WorkerAnswerHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request, const userver::formats::json::Value&,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kGet: {
        const std::string& requestId = request.GetArg("requestId");

        auto workerAnswer = m_taskScheduler.getAnswer(requestId);
        if (!workerAnswer) {
            request.SetResponseStatus(userver::server::http::HttpStatus::kNotFound);
            return {};
        }

        request.SetResponseStatus(userver::server::http::HttpStatus::kOk);
        return userver::formats::json::ValueBuilder(workerAnswer.value()).ExtractValue();
    }
    default: {
        throw userver::server::handlers::ClientError(
            userver::server::handlers::ExternalBody { fmt::format("Unsupported method {}", request.GetMethod()) });
    }
    }
}

} // namespace Manager
