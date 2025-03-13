#include "services/lexicographic_word_generator.hpp"

#include <cmath>

namespace Worker {

LexicographicWordGenerator::LexicographicWordGenerator(const std::string& alphabet, size_t maxLength)
    : m_alphabet(alphabet)
    , maxLength_(maxLength)
    , totalWords_(getTotalWordsCount())
{
}

std::string LexicographicWordGenerator::getWordAtIndex(size_t index) const
{
    if (index >= totalWords_) {
        return "";
    }
    return generateWordFromIndex(index);
}

std::string LexicographicWordGenerator::generateWordFromIndex(size_t index) const
{
    std::string word;
    size_t alphabetSize = m_alphabet.size();

    while (index > 0 || word.size() < maxLength_) {
        size_t letterIndex = index % alphabetSize;
        word = m_alphabet[letterIndex] + word;
        index /= alphabetSize;
    }

    while (word.size() < maxLength_) {
        word = m_alphabet[0] + word;
    }

    return word;
}

size_t LexicographicWordGenerator::getTotalWordsCount() const
{
    size_t total = 0;
    for (size_t length = 1; length <= maxLength_; ++length) {
        total += std::pow(m_alphabet.size(), length);
    }
    return total;
}

} // namespace Worker
