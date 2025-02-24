#include "models/worker.hpp"

namespace Manager {

Worker::Worker(std::shared_ptr<WorkerConnection> connection)
    : m_connection(connection)
{
}

void Worker::Process(const Task& task)
{
    m_connection->Send(task);
}

} // Manager
