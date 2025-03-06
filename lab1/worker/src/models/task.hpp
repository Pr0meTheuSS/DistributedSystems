#pragma once

#include <string>

#include <userver/formats/json/value_builder.hpp>

#include "models/task.hpp"

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Worker {

struct Task final {
    std::string requestId;
    std::string hash;
    std::size_t maxLength;
};

userver::formats::json::Value Serialize(const Task& task,
    userver::formats::serialize::To<userver::formats::json::Value>);

} // Manager
