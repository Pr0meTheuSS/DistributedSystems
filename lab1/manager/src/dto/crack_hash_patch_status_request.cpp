#include <dto/crack_hash_patch_status_request.hpp>

#include <userver/formats/json/value_builder.hpp>

namespace {

// std::string CrackStatusToString(const CrackStatus& status)
// {
//     switch (status) {
//     case CrackStatus::IN_PROGESS:
//         return "IN_PROGRESS";
//     case CrackStatus::ERROR:
//         return "ERROR";
//     case CrackStatus::READY:
//         return "READY";
//     case CrackStatus::UNDEFINED:
//         return "UNDEFINED";

//     default:
//         return {};
//     }
// }

Manager::CrackStatus CrackStatusFromString(const std::string& statusString)
{
    if ("IN_PROGRESS" == statusString) {
        return Manager::CrackStatus::IN_PROGRESS;

    } else if ("READY" == statusString) {
        return Manager::CrackStatus::READY;
    } else if ("ERROR" == statusString) {
        return Manager::CrackStatus::ERROR;
    }
    return Manager::CrackStatus::UNDEFINED;
}
}

namespace Manager {

CrackHashPatchStatusRequest Parse(const userver::formats::json::Value& json,
    userver::formats::parse::To<CrackHashPatchStatusRequest>)
{
    if (!json.HasMember("requestId") || !json.HasMember("status")) {
        return {};
    }

    return CrackHashPatchStatusRequest {
        json["requestId"].As<std::string>(),
        CrackStatusFromString(json["status"].As<std::string>()),
    };
}

} // Manager
