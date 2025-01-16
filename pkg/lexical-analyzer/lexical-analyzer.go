package lexicalanalyzer

import (
	"fmt"
	"regexp"
	"strings"

	sourcetext "gitgub.com/aswait/go-syntactic-analyzer/pkg/source-text"
)

type LexicalAnalyzerer interface {
	Transliterate() (string, error)
	Validate() string
	SourceLoadFromInput(input string)
	Analyze() ([]Token, error)
}

type LexicalAnalyzer struct {
	source          sourcetext.SourceTexter
	alphabet        map[rune]string
	transliterated  []rune
	currentPosition int
	reservedWords   map[string]string
}

type Token struct {
	Value string
	Type  string
}

func NewLexicalAnalyzer(source sourcetext.SourceTexter) *LexicalAnalyzer {
	alphabet := make(map[rune]string)

	for r := 'a'; r <= 'z'; r++ {
		alphabet[r] = "Letter"
	}
	for r := 'A'; r <= 'Z'; r++ {
		alphabet[r] = "Letter"
	}

	for r := '0'; r <= '9'; r++ {
		alphabet[r] = "Digit"
	}

	alphabet[' '] = "Space"
	alphabet['\n'] = "EndRow"
	alphabet['\t'] = "Tab"
	alphabet[';'] = "Semicolon"
	alphabet[','] = "Comma"
	alphabet['='] = "Assign"
	alphabet['+'] = "Plus"
	alphabet['*'] = "Multiply"
	alphabet['('] = "LeftParen"
	alphabet[')'] = "RightParen"
	alphabet['/'] = "Comment"
	alphabet['*'] = "Comment"

	reservedWords := map[string]string{
		"let": "Let",
		";":   "Semicolon",
		",":   "Comma",
		"=":   "Assign",
		"+":   "Plus",
		"*":   "Multiply",
		"(":   "LeftParen",
		")":   "RightParen",
	}

	return &LexicalAnalyzer{
		source:          source,
		alphabet:        alphabet,
		transliterated:  []rune{},
		reservedWords:   reservedWords,
		currentPosition: 0,
	}
}

func (la *LexicalAnalyzer) Transliterate() (string, error) {
	var result strings.Builder

	for la.source.HasMoreSymbols() {
		symbol, err := la.source.ReadNextSymbol()
		if err != nil {
			return "", fmt.Errorf("Error reading text: %v", err)
		}

		class, exists := la.alphabet[symbol]
		if !exists {
			return "", fmt.Errorf("Символ '%c' не принадлежит алфавиту", symbol)
		}
		la.transliterated = append(la.transliterated, symbol)
		if symbol == '\n' || symbol == '\t' {
			result.WriteString(fmt.Sprintf("(%s)\n", class))
		} else {
			result.WriteString(fmt.Sprintf("(%s, %c)\n", class, symbol))
		}

	}

	return result.String(), nil
}

func (la *LexicalAnalyzer) Validate() string {
	if la.source.HasMoreSymbols() {
		return "Error: Text not fully processed."
	}
	return "Текст верен"
}

func (la *LexicalAnalyzer) SourceLoadFromInput(input string) {
	la.source.LoadFromInput(input)
	la.transliterated = []rune{}
	la.currentPosition = 0
}

func (la *LexicalAnalyzer) Analyze() ([]Token, error) {
	var tokens []Token
	commentMode := false
	for la.currentPosition < len(la.transliterated) {
		symbol := la.transliterated[la.currentPosition]
		la.currentPosition++
		if commentMode {
			if symbol == '*' && la.currentPosition < len(la.transliterated) && la.transliterated[la.currentPosition] == '/' {
				commentMode = false
				la.currentPosition++
			}
			continue
		}
		if symbol == '/' && la.currentPosition < len(la.transliterated) && la.transliterated[la.currentPosition] == '*' {
			commentMode = true
			la.currentPosition++
			continue
		}
		if symbol == ' ' || symbol == '\n' || symbol == '\t' {
			continue
		}
		word := string(symbol)
		for la.currentPosition < len(la.transliterated) {
			nextSymbol := la.transliterated[la.currentPosition]
			if nextSymbol == ' ' || nextSymbol == '\n' || nextSymbol == '\t' || nextSymbol == ';' || nextSymbol == ',' {
				break
			}
			word += string(nextSymbol)
			la.currentPosition++
		}

		if matched, _ := regexp.MatchString(`^(010)*100(000)*$`, word); matched {
			tokens = append(tokens, Token{Value: word, Type: "Number"})
		} else if matched, _ := regexp.MatchString(`^cd[a-d]*$`, word); matched {
			tokens = append(tokens, Token{Value: word, Type: "Identifier"})
		} else if t, exists := la.reservedWords[word]; exists {
			tokens = append(tokens, Token{Value: word, Type: t})
		} else {
			return nil, fmt.Errorf("Несоответствие слова: %s", word)
		}
	}
	return tokens, nil
}
