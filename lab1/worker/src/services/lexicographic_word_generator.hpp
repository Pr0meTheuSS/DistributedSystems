#pragma once

#include <string>

namespace Worker {

class LexicographicWordGenerator {
public:
    LexicographicWordGenerator(const std::string& alphabet, size_t maxLength);

    std::string getWordAtIndex(size_t index) const;
    size_t getTotalWordsCount() const;

private:
    std::string generateWordFromIndex(size_t index) const;

    std::string m_alphabet;
    size_t maxLength_;
    size_t totalWords_;
};

} // namespace Worker
