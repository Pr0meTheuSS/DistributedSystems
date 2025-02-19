#include <models/task.hpp>

#include <string>

#include <userver/formats/json/value_builder.hpp>

#include <models/task.hpp>

namespace Manager {

userver::formats::json::Value Serialize(const Task& task,
    userver::formats::serialize::To<userver::formats::json::Value>)
{
    userver::formats::json::ValueBuilder builder;

    builder["hash"] = task.hash;
    builder["maxLength"] = task.maxLength;
    builder["requestId"] = task.requestId;
    builder["offset"] = task.offset;
    builder["size"] = task.size;

    return builder.ExtractValue();
}

} // Manager
