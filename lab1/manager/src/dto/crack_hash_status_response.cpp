#include "dto/crack_hash_status_response.hpp"

#include <userver/formats/json/value_builder.hpp>

namespace {

std::string CrackStatusToString(const Manager::CrackStatus& status)
{
    switch (status) {
    case Manager::CrackStatus::IN_PROGRESS:
        return "IN_PROGRESS";
    case Manager::CrackStatus::ERROR:
        return "ERROR";
    case Manager::CrackStatus::READY:
        return "READY";
    case Manager::CrackStatus::UNDEFINED:
        return "UNDEFINED";

    default:
        return {};
    }
}

}

namespace Manager {

userver::formats::json::Value Serialize(
    const CrackHashStatusResponse& response, userver::formats::serialize::To<userver::formats::json::Value>)
{
    userver::formats::json::ValueBuilder builder;

    builder["requestId"] = response.requestId;
    builder["status"] = CrackStatusToString(response.status);

    return builder.ExtractValue();
}

} // namespace Manager

