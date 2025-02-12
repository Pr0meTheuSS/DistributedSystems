#include "hash_crack_handler.hpp"

#include <iostream>

#include <userver/server/handlers/exceptions.hpp>

#include "dto/crack_hash_request.hpp"

namespace samples::hello {

userver::formats::json::Value CrackHashHandler::HandleRequestJsonThrow(
    const userver::server::http::HttpRequest& request,
    const userver::formats::json::Value& requestBody,
    userver::server::request::RequestContext&) const
{
    switch (request.GetMethod()) {
    case userver::server::http::HttpMethod::kPost: {
        // const std::string& name = request.GetArg("name");
        // if (name.empty()) {
        //     throw userver::server::handlers::ClientError(
        //         userver::server::handlers::ExternalBody{ .body = "No parameter 'name' in path arg"});
        // }

        auto crackHashRequest = requestBody.As<CrackHashRequest>();

        request.SetResponseStatus(userver::server::http::HttpStatus::kCreated);
        return userver::formats::json::ValueBuilder(crackHashRequest).ExtractValue();
    }
    default: {
        throw userver::server::handlers::ClientError(userver::server::handlers::ExternalBody {
            .body = fmt::format("Unsupported method {}", request.GetMethod()) });
    }
    }
}

} // namespace samples::hello