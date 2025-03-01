#include "workers_factory_component.hpp"

#include <memory>
#include <string>
#include <vector>

#include <libenvpp/env.hpp>

#include "models/worker_connection.hpp"

namespace Manager {

WorkersFactoryComponent::WorkersFactoryComponent(
    const userver::components::ComponentConfig& config,
    const userver::components::ComponentContext& context)
    : userver::components::LoggableComponentBase(config, context)
{
    // Получаем количество рабочих из переменных окружения
    const auto workersAmount = env::get_or<unsigned int>("WORKERS_AMOUNT", 0u);

    for (auto i = 1u; i <= workersAmount; i++) {
        // Получаем URL для каждого рабочего
        std::string workerUrl = env::get<std::string>("WORKER_URL_" + std::to_string(i)).value();

        // Создаем подключение для рабочего
        auto connection = std::make_unique<HttpWorkerConnection>(config, context);
        connection->SetUrl(workerUrl);

        // Сохраняем подключение в список рабочих
        m_workers.push_back(Worker(std::move(connection))); // Передаем владение в Worker
    }
}

std::vector<Worker> WorkersFactoryComponent::getWorkers() const
{
    return m_workers;
}

} // namespace Manager
