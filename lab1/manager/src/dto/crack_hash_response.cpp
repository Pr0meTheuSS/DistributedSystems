#include "dto/crack_hash_response.hpp"

#include <userver/formats/json/value_builder.hpp>

namespace Manager {

userver::formats::json::Value Serialize(const CrackHashResponse& response,
    userver::formats::serialize::To<userver::formats::json::Value>)
{
    userver::formats::json::ValueBuilder builder;
    builder["requestId"] = response.requestId;
    return builder.ExtractValue();
}

} // namespace Manager

