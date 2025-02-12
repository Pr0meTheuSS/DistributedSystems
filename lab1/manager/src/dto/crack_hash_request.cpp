#include "dto/crack_hash_request.hpp"

#include <userver/formats/json/value_builder.hpp>

namespace samples::hello {

CrackHashRequest Parse(const userver::formats::json::Value& json,
    userver::formats::parse::To<CrackHashRequest>)
{
    if (!json.HasMember("hash") || !json.HasMember("maxLength")) {
        return {};
    }

    return CrackHashRequest {
        .hash = json["hash"].As<std::string>(),
        .maxLength = json["maxLength"].As<int>()
    };
}

userver::formats::json::Value Serialize(const CrackHashRequest& data,
    userver::formats::serialize::To<userver::formats::json::Value>)
{
    userver::formats::json::ValueBuilder builder;

    builder["hash"] = data.hash;
    builder["maxLength"] = data.maxLength;

    return builder.ExtractValue();
}

} // namespace samples::hello
