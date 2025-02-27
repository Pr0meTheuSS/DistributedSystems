#pragma once

#include <string>
#include <vector>

#include <userver/formats/parse/to.hpp>
#include <userver/formats/serialize/to.hpp>

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Manager {

struct WorkerAnswer final {
    std::string requestId;
    std::size_t partNumber;
    std::vector<std::string> words;
};

WorkerAnswer Parse(const userver::formats::json::Value&,
    userver::formats::parse::To<WorkerAnswer>);

userver::formats::json::Value Serialize(const WorkerAnswer&,
    userver::formats::serialize::To<userver::formats::json::Value>);

} // namespace Manager
