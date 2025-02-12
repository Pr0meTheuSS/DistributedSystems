#pragma once

#include <string>

#include <userver/formats/parse/to.hpp>
#include <userver/formats/serialize/to.hpp>

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace samples::hello {

struct CrackHashRequest final {
    std::string hash;
    int maxLength;
};

CrackHashRequest Parse(const userver::formats::json::Value&,
    userver::formats::parse::To<CrackHashRequest>);

userver::formats::json::Value Serialize(const CrackHashRequest&,
    userver::formats::serialize::To<userver::formats::json::Value>);
}
