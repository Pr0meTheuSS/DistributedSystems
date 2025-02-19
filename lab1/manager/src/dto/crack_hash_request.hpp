#pragma once

#include <string>

#include <userver/formats/parse/to.hpp>

namespace userver::formats::json {
class Value;
} // namespace userver::formats::json

namespace Manager {

struct CrackHashRequest final {
    std::string hash;
    int maxLength;
};

CrackHashRequest Parse(const userver::formats::json::Value&,
    userver::formats::parse::To<CrackHashRequest>);

}
