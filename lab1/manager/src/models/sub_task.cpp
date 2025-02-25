#include "models/task.hpp"

#include <string>

#include <userver/formats/json/value_builder.hpp>

#include "models/sub_task.hpp"

namespace Manager {

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
