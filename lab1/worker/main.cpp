#include <userver/clients/dns/component.hpp>
#include <userver/components/minimal_server_component_list.hpp>
#include <userver/utest/using_namespace_userver.hpp>
#include <userver/utils/daemon_run.hpp>

// Note: this is for the purposes of tests/samples only
#include <userver/utest/using_namespace_userver.hpp>

#include "clients/http_client.hpp"
#include "components/http_worker_connection.hpp"
#include "handlers/crack_hash_handler.hpp"

int main(int argc, char* argv[])
{
    auto component_list = components::MinimalServerComponentList()
                              .Append<userver::components::HttpClient>()
                              .Append<userver::clients::dns::Component>()
                              .Append<Worker::HttpClient>()
                              .Append<Worker::CrackHashHandler>();

    return utils::DaemonMain(argc, argv, component_list);
}
