#include "handlers/crack_hash_status_handler.hpp"

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include "dto/crack_hash_patch_status_request.hpp"
#include "dto/worker_answer.hpp"

namespace Manager {

CrackHashStatusHandler::CrackHashStatusHandler(const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::server::handlers::HttpHandlerJsonBase(config, context)
    , m_taskScheduler(context.FindComponent<TaskSchedulerComponent>()) { };

userver::formats::json::Value CrackHashStatusHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request,
    const userver::formats::json::Value& requestBody,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kPatch: {
        auto req = requestBody.As<CrackHashPatchStatusRequest>();

        CrackRequest requestForUpdate;
        requestForUpdate.id = req.id;
        requestForUpdate.status = req.status;
        if (!m_taskScheduler.UpdateStatus(requestForUpdate)) {
            throw userver::server::handlers::ClientError(userver::server::handlers::ExternalBody {
                fmt::format("Cannot update status") });
        }
        request.SetResponseStatus(userver::server::http::HttpStatus::kNoContent);
        return {};
    }
    case userver::server::http::HttpMethod::kPost: {
        WorkerAnswer workerAnswer = requestBody.As<WorkerAnswer>();

        m_taskScheduler.completeSubTask(workerAnswer.requestId, workerAnswer.partNumber);
        m_taskScheduler.saveWorkerAnswer(workerAnswer);

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
