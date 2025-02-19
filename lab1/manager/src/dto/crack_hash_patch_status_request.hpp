#pragma once

#include <string>

#include <userver/formats/parse/to.hpp>

#include <models/crack_request.hpp>

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Manager {

struct CrackHashPatchStatusRequest final {
    std::string id;
    CrackStatus status;
};

CrackHashPatchStatusRequest Parse(const userver::formats::json::Value&,
    userver::formats::parse::To<CrackHashPatchStatusRequest>);

} // Manager

