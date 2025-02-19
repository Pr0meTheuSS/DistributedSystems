#pragma once

#include <string>

#include <userver/formats/json/value_builder.hpp>

#include <models/task.hpp>

namespace Manager {

struct Task final {
    std::string hash;
    std::string requestId;
    std::size_t maxLength;
    std::size_t offset;
    std::size_t size;
};

userver::formats::json::Value Serialize(const Task& task,
    userver::formats::serialize::To<userver::formats::json::Value>);

} // Manager
