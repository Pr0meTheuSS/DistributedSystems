#include <iomanip>
#include <iostream>
#include <openssl/md5.h> // Библиотека для вычисления MD5
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>

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

    // Метод для получения i-го слова в лексикографическом порядке
    std::string getWordAtIndex(size_t index) const
    {
        if (index >= totalWords_) {
            return ""; // Если индекс выходит за пределы, возвращаем пустую строку
        }
        return generateWordFromIndex(index);
    }

    // Генерация слова по индексу
    std::string generateWordFromIndex(size_t index) const
    {
        std::string word;
        size_t alphabetSize = alphabet_.size();

        // Генерация каждого символа слова
        while (index > 0 || word.size() < maxLength_) {
            size_t letterIndex = index % alphabetSize;
            word = alphabet_[letterIndex] + word;
            index /= alphabetSize;
        }

        // Заполнение недостающих символов в начале слова (если оно короче maxLength)
        while (word.size() < maxLength_) {
            word = alphabet_[0] + word;
        }

        return word;
    }

    // Получение общего числа возможных слов
    size_t getTotalWordsCount() const
    {
        size_t total = 0;
        for (size_t length = 1; length <= maxLength_; ++length) {
            total += std::pow(alphabet_.size(), length);
        }
        return total;
    }

private:
    std::string alphabet_; // Алфавит для генерации слов
    size_t maxLength_; // Максимальная длина слов
    size_t totalWords_; // Общее количество возможных слов
};

class WorkerService {
public:
    WorkerService(const std::string& alphabet)
        : alphabet_(alphabet)
        , generator_(alphabet, 0)
    {
    }

    void process(const SubTask& subTask)
    {
        generator_ = LexicographicWordGenerator { alphabet_, subTask.maxLength };

        std::unordered_map<std::string, WorkerAnswer> resultMap;

        size_t totalWords = generator_.getTotalWordsCount();
        std::cout << "Total Words: " << totalWords << std::endl;

        size_t wordsPerPart = totalWords / subTask.partCount;
        size_t remainingWords = totalWords % subTask.partCount;

        std::cout << "subTask.partNumber: " << subTask.partNumber << std::endl;
        std::cout << "wordsPerPart: " << wordsPerPart << ", remainingWords: " << remainingWords << std::endl;

        // Рассчитываем startIndex и endIndex для текущей части
        size_t startIndex = wordsPerPart * subTask.partNumber + std::min(subTask.partNumber, remainingWords);
        size_t endIndex = startIndex + wordsPerPart + (subTask.partNumber < remainingWords ? 1 : 0);

        std::cout << "start and end: " << startIndex << " : " << endIndex << std::endl;

        std::vector<std::string> words;
        for (size_t i = startIndex; i < endIndex; ++i) {
            std::string word = generator_.getWordAtIndex(i);
            std::cout << "Word #" << i << ": " << word << std::endl;

            if (word.empty())
                break;

            // Проверка на совпадение с хешем
            if (computeMD5(word) == subTask.hash) {
                words.push_back(word);
            }
        }

        WorkerAnswer answer = { subTask.requestId, subTask.partNumber, words };
        resultMap[subTask.requestId] = answer;

        for (const auto& entry : resultMap) {
            std::cout << "Request ID: " << entry.second.requestId
                      << ", Part Number: " << entry.second.partNumber << "\nWords: ";
            for (const auto& word : entry.second.words) {
                std::cout << word << " ";
            }
            std::cout << std::endl;
        }
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

// int main() {
//     // Пример данных для SubTask
//     SubTask subTask = {
//         "task1",       // requestId
//         "9e107d9d372bb6826bd81d3542a419d6", // хэш MD5
//         3,             // maxLength
//         2,             // partCount
//         1              // partNumber
//     };

//     WorkerService workerService("abcdefghijklmnopqrstuvwxyz0123456789");
//     workerService.process(subTask);

//     return 0;
// }

} // namespace Worker
