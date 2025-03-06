#pragma once

#include <string>

#include <userver/formats/json/value_builder.hpp>

#include "models/task.hpp"

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Worker {

struct SubTask final {
    std::string requestId;
    std::string hash;

    std::size_t maxLength;
    std::size_t partCount;
    std::size_t partNumber;
};

userver::formats::json::Value Serialize(const SubTask&,
    userver::formats::serialize::To<userver::formats::json::Value>);

SubTask Parse(const userver::formats::json::Value&,
    userver::formats::parse::To<SubTask>);

} // Manager
