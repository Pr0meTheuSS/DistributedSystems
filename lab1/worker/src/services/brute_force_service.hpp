#pragma once

#include <ostream>
#include <string>
#include <unordered_map>
#include <vector>

#include "dto/worker_answer.hpp"
#include "models/sub_task.hpp"
#include "services/lexicographic_word_generator.hpp"

namespace Worker {

class BruteForceService {
public:
    BruteForceService(const std::string& alphabet, std::ostream& logStream);
    virtual WorkerAnswer process(const SubTask& subTask);

private:
    std::string computeMD5(const std::string& str);

    std::ostream& m_logStream;
    std::string m_alphabet;
    LexicographicWordGenerator m_generator;
};

} // namespace Worker
