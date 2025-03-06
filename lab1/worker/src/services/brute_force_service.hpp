#pragma once
#include <iomanip>
#include <iostream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>

#include <openssl/md5.h>

#include "dto/worker_answer.hpp"
#include "models/sub_task.hpp"

namespace Worker {
class LexicographicWordGenerator {
public:
    LexicographicWordGenerator(const std::string& alphabet, size_t maxLength)
        : alphabet_(alphabet)
        , maxLength_(maxLength)
    {
        totalWords_ = getTotalWordsCount();
    }

    std::string getWordAtIndex(size_t index) const
    {
        if (index >= totalWords_) {
            return "";
        }
        return generateWordFromIndex(index);
    }

    std::string generateWordFromIndex(size_t index) const
    {
        std::string word;
        size_t alphabetSize = alphabet_.size();

        while (index > 0 || word.size() < maxLength_) {
            size_t letterIndex = index % alphabetSize;
            word = alphabet_[letterIndex] + word;
            index /= alphabetSize;
        }

        while (word.size() < maxLength_) {
            word = alphabet_[0] + word;
        }

        return word;
    }

    size_t getTotalWordsCount() const
    {
        size_t total = 0;
        for (size_t length = 1; length <= maxLength_; ++length) {
            total += std::pow(alphabet_.size(), length);
        }
        return total;
    }

private:
    std::string alphabet_;
    size_t maxLength_;
    size_t totalWords_;
};

class BruteForceService {
public:
    BruteForceService(const std::string& alphabet)
        : alphabet_(alphabet)
        , generator_(alphabet, 0)
    {
    }

    virtual WorkerAnswer process(const SubTask& subTask)
    {
        generator_ = LexicographicWordGenerator { alphabet_, subTask.maxLength };

        std::unordered_map<std::string, WorkerAnswer> resultMap;

        size_t totalWords = generator_.getTotalWordsCount();
        std::cout << "Total Words: " << totalWords << std::endl;

        size_t wordsPerPart = totalWords / subTask.partCount;
        size_t remainingWords = totalWords % subTask.partCount;

        std::cout << "subTask.partNumber: " << subTask.partNumber << std::endl;
        std::cout << "wordsPerPart: " << wordsPerPart << ", remainingWords: " << remainingWords << std::endl;

        size_t startIndex = wordsPerPart * subTask.partNumber + std::min(subTask.partNumber, remainingWords);
        size_t endIndex = startIndex + wordsPerPart + (subTask.partNumber < remainingWords ? 1 : 0);

        std::cout << "start and end: " << startIndex << " : " << endIndex << std::endl;

        std::vector<std::string> words;
        for (size_t i = startIndex; i < endIndex; ++i) {
            std::string word = generator_.getWordAtIndex(i);
            std::cout << "Word #" << i << ": " << word << std::endl;

            if (word.empty())
                break;

            if (computeMD5(word) == subTask.hash) {
                words.push_back(word);
            }
        }

        WorkerAnswer answer = { subTask.requestId, subTask.partNumber, words };
        return resultMap[subTask.requestId] = answer;
    }

private:
    std::string alphabet_;
    LexicographicWordGenerator generator_;

    std::string computeMD5(const std::string& str)
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
};

} // namespace Worker
