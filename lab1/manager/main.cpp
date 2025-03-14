#include <userver/clients/dns/component.hpp>
#include <userver/components/minimal_server_component_list.hpp>
#include <userver/utest/using_namespace_userver.hpp>
#include <userver/utils/daemon_run.hpp>

#include <userver/utest/using_namespace_userver.hpp>

#include "clients/http_client.hpp"
#include "components/http_worker_connection.hpp"
#include "components/task_scheduler_component.hpp"
#include "components/workers_factory_component.hpp"
#include "handlers/crack_hash_handler.hpp"
#include "handlers/crack_hash_status_handler.hpp"
#include "handlers/hash_status_handler.hpp"
#include "handlers/worker_answer_handler.hpp"
#include "repositories/requests_repository_in_memory.hpp"
#include "repositories/tasks_repository_in_memory.hpp"

int main(int argc, char* argv[])
{
    auto component_list = components::MinimalServerComponentList()
                              .Append<Manager::TasksRepositoryInMemory>()
                              .Append<Manager::WorkersFactoryComponent>()
                              .Append<userver::components::HttpClient>()
                              .Append<userver::clients::dns::Component>()
                              .Append<Manager::HttpClient>()
                              .Append<Manager::HttpWorkerConnection>()
                              .Append<Manager::WorkerAnswerHandler>()
                              .Append<Manager::CrackHashHandler>()
                              .Append<Manager::CrackHashStatusHandler>()
                              .Append<Manager::HashStatusHandler>()
                              .Append<Manager::TaskSchedulerComponent>()
                              .Append<Manager::RequestsRepositoryInMemory>();

    return utils::DaemonMain(argc, argv, component_list);
}
