#include "models/task.hpp"

#include <string>

#include <userver/formats/json/value_builder.hpp>

#include "models/sub_task.hpp"

namespace Worker {

SubTask Parse(const userver::formats::json::Value& json,
    userver::formats::parse::To<SubTask>)
{
    if (!json.HasMember("requestId")
        || !json.HasMember("hash")
        || !json.HasMember("maxLength")
        || !json.HasMember("partNumber")
        || !json.HasMember("partCount")) {
        return {};
    }

    return SubTask {
        json["requestId"].As<std::string>(),
        json["hash"].As<std::string>(),
        json["maxLength"].As<std::size_t>(),
        json["partCount"].As<std::size_t>(),
        json["partNumber"].As<std::size_t>()
    };
}

userver::formats::json::Value Serialize(const SubTask& task,
    userver::formats::serialize::To<userver::formats::json::Value>)
{
    userver::formats::json::ValueBuilder builder;

    builder["hash"] = task.hash;
    builder["maxLength"] = task.maxLength;
    builder["requestId"] = task.requestId;
    builder["partNumber"] = task.partNumber;
    builder["partCount"] = task.partCount;

    return builder.ExtractValue();
}

} // Manager
