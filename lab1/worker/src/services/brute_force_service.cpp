#include "services/brute_force_service.hpp"

#include <iomanip>
#include <iostream>
#include <openssl/md5.h>
#include <sstream>

namespace Worker {

BruteForceService::BruteForceService(const std::string& alphabet, std::ostream& logStream)
    : m_alphabet(alphabet)
    , m_generator(alphabet, 0)
    , m_logStream(logStream)
{
}

WorkerAnswer BruteForceService::process(const SubTask& subTask)
{
    m_generator = LexicographicWordGenerator { m_alphabet, subTask.maxLength };

    std::unordered_map<std::string, WorkerAnswer> resultMap;
    size_t totalWords = m_generator.getTotalWordsCount();

    m_logStream << "Total Words: " << totalWords << std::endl;

    size_t wordsPerPart = totalWords / subTask.partCount;
    size_t remainingWords = totalWords % subTask.partCount;

    m_logStream << "subTask.partNumber: " << subTask.partNumber << std::endl;
    m_logStream << "wordsPerPart: " << wordsPerPart << ", remainingWords: " << remainingWords << std::endl;

    size_t startIndex = wordsPerPart * subTask.partNumber + std::min(subTask.partNumber, remainingWords);
    size_t endIndex = startIndex + wordsPerPart + (subTask.partNumber < remainingWords ? 1 : 0);

    m_logStream << "start and end: " << startIndex << " : " << endIndex << std::endl;

    std::vector<std::string> words;
    for (size_t i = startIndex; i < endIndex; ++i) {
        std::string word = m_generator.getWordAtIndex(i);
        m_logStream << "Word #" << i << ": " << word << std::endl;

        if (word.empty())
            break;

        if (computeMD5(word) == subTask.hash) {
            words.push_back(word);
        }
    }

    WorkerAnswer answer = { subTask.requestId, subTask.partNumber, words };
    return resultMap[subTask.requestId] = answer;
}

std::string BruteForceService::computeMD5(const std::string& str)
{
    unsigned char result[MD5_DIGEST_LENGTH];
    MD5_CTX md5Context;
    MD5_Init(&md5Context);
    MD5_Update(&md5Context, str.c_str(), str.length());
    MD5_Final(result, &md5Context);

    std::stringstream ss;
    for (int i = 0; i < MD5_DIGEST_LENGTH; ++i) {
        ss << std::setw(2) << std::setfill('0') << std::hex << (int)result[i];
    }
    return ss.str();
}

} // namespace Worker
