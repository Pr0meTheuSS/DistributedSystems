#include "models/worker.hpp"

namespace Manager {

Worker::Worker(std::shared_ptr<WorkerConnection> connection)
    : m_connection(connection)
    , m_isReady(true)
{
}

void Worker::Process(const SubTask& task)
{
    m_connection->Send(task);
}

bool Worker::isReady() const
{
    return m_isReady;
}

void Worker::setReady()
{
    m_isReady = true;
}

} // Manager
