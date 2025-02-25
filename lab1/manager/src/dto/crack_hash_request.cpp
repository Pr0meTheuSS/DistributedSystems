#include "dto/crack_hash_request.hpp"

#include <userver/formats/json/value_builder.hpp>

namespace Manager {

CrackHashRequest Parse(const userver::formats::json::Value& json,
    userver::formats::parse::To<CrackHashRequest>)
{
    if (!json.HasMember("hash") || !json.HasMember("maxLength")) {
        return {};
    }

    return CrackHashRequest {
        json["hash"].As<std::string>(),
        json["maxLength"].As<std::size_t>()
    };
}

} // namespace samples::hello
