#pragma once

#include <string>

#include <userver/formats/serialize/to.hpp>

#include "models/crack_request.hpp"

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Manager {

struct CrackHashStatusResponse final {
    std::string requestId;
    CrackStatus status;
};

userver::formats::json::Value Serialize(
    const CrackHashStatusResponse&, userver::formats::serialize::To<userver::formats::json::Value>);
}
