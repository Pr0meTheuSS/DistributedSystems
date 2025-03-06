#pragma once

#include <string>

#include <userver/formats/serialize/to.hpp>

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Manager {

struct CrackHashResponse final {
    std::string requestId;
};

userver::formats::json::Value Serialize(
    const CrackHashResponse&, userver::formats::serialize::To<userver::formats::json::Value>);
}
