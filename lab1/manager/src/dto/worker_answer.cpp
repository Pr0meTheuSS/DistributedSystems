#include "dto/worker_answer.hpp"

#include <userver/formats/json/value_builder.hpp>
#include <userver/formats/parse/common_containers.hpp>
#include <userver/formats/serialize/common_containers.hpp>

namespace Manager {

WorkerAnswer Parse(const userver::formats::json::Value& json, userver::formats::parse::To<WorkerAnswer>)
{
    if (!json.HasMember("requestId") || !json.HasMember("partNumber") || !json.HasMember("words")) {
        return {};
    }

    return WorkerAnswer { json["requestId"].As<std::string>(), json["partNumber"].As<std::size_t>(),
        json["words"].As<std::vector<std::string>>() };
}

userver::formats::json::Value Serialize(
    const WorkerAnswer& answer, userver::formats::serialize::To<userver::formats::json::Value>)
{
    userver::formats::json::ValueBuilder builder;

    builder["requestId"] = answer.requestId;
    builder["partNumber"] = answer.partNumber;
    builder["words"] = answer.words;

    return builder.ExtractValue();
}

} // namespace Manager
