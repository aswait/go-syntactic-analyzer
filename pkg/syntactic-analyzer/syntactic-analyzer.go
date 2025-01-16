package syntacticanalyzer

import (
	"fmt"

	lexicalanalyzer "gitgub.com/aswait/go-syntactic-analyzer/pkg/lexical-analyzer"
)

type SyntacticAnalyzer struct {
	tokens       []lexicalanalyzer.Token
	currentIndex int
	initialized  map[string]bool
}

func NewSyntacticAnalyzer() *SyntacticAnalyzer {
	return &SyntacticAnalyzer{
		initialized: make(map[string]bool),
	}
}

func (sa *SyntacticAnalyzer) currentToken() lexicalanalyzer.Token {
	if sa.currentIndex < len(sa.tokens) {
		fmt.Printf("Текущий токен: %v\n", sa.tokens[sa.currentIndex])
		return sa.tokens[sa.currentIndex]
	}
	fmt.Println("Достигнут EOF")
	return lexicalanalyzer.Token{Type: "EOF"}
}

func (sa *SyntacticAnalyzer) nextToken() {
	sa.currentIndex++
	if sa.currentIndex < len(sa.tokens) {
		fmt.Printf("Переход к следующему токену: %v\n", sa.tokens[sa.currentIndex])
	} else {
		fmt.Println("Достигнут конец токенов")
	}
}

func (sa *SyntacticAnalyzer) Parse(tokens []lexicalanalyzer.Token) error {
	sa.tokens = tokens
	sa.currentIndex = 0

	fmt.Println("Входные токены:", tokens)

	if err := sa.parseS(); err != nil {
		return fmt.Errorf("Ошибка синтаксического анализа: %v", err)
	}

	if sa.currentIndex < len(sa.tokens) {
		return fmt.Errorf("неожиданный токен '%s' на позиции %d", sa.tokens[sa.currentIndex].Value, sa.currentIndex)
	}

	fmt.Println("Парсинг успешно завершён")
	return nil
}

func (sa *SyntacticAnalyzer) parseS() error {
	if err := sa.parseD(); err != nil {
		return err
	}
	return sa.parseO()
}

func (sa *SyntacticAnalyzer) parseD() error {
	token := sa.currentToken()
	if token.Type != "Let" {
		return fmt.Errorf("ожидалось 'let', но найдено '%s'", token.Value)
	}
	fmt.Printf("Обнаружено ключевое слово: %s\n", token.Value)
	sa.nextToken()

	if err := sa.parseL(); err != nil {
		return err
	}

	token = sa.currentToken()
	if token.Type != "Semicolon" {
		return fmt.Errorf("ожидался ';', но найдено '%s'", token.Value)
	}
	fmt.Printf("Обнаружен конец инструкции: %s\n", token.Value)
	sa.nextToken()
	return nil
}

func (sa *SyntacticAnalyzer) parseL() error {
	fmt.Println("Начало анализа списка идентификаторов")
	if err := sa.parseV(); err != nil {
		return err
	}

	for {
		token := sa.currentToken()
		if token.Type == "Comma" {
			fmt.Printf("Обнаружен разделитель списка: %s\n", token.Value)
			sa.nextToken()
			if err := sa.parseV(); err != nil {
				return err
			}
		} else {
			break
		}
	}

	fmt.Println("Конец анализа списка идентификаторов")
	return nil
}

func (sa *SyntacticAnalyzer) parseV() error {
	token := sa.currentToken()
	fmt.Printf("Анализ идентификатора/числа: %v\n", token)

	if token.Type != "Identifier" {
		return fmt.Errorf("ожидался идентификатор, но найдено '%s'", token.Value)
	}

	varName := token.Value
	sa.nextToken()

	token = sa.currentToken()
	if token.Type == "Assign" {
		fmt.Printf("Обнаружен оператор присваивания: %s\n", token.Value)
		sa.nextToken()

		if err := sa.parseE(); err != nil {
			return err
		}

		sa.initialized[varName] = true
	} else if !sa.initialized[varName] {
		return fmt.Errorf("переменная '%s' не проинициализирована", varName)
	}

	return nil
}

func (sa *SyntacticAnalyzer) parseO() error {
	if sa.currentToken().Type == "EOF" {
		return nil // Завершаем обработку, если достигнут конец токенов
	}
	if err := sa.parseP(); err != nil {
		return err
	}
	if sa.currentToken().Type == "Semicolon" {
		sa.nextToken()
		return sa.parseO()
	}
	return nil
}

func (sa *SyntacticAnalyzer) parseP() error {
	token := sa.currentToken()
	if token.Type != "Identifier" {
		return fmt.Errorf("ожидался идентификатор, но найдено '%s'", token.Value)
	}
	varName := token.Value
	sa.nextToken()
	token = sa.currentToken()
	if token.Type != "Assign" {
		return fmt.Errorf("ожидался '=', но найдено '%s'", token.Value)
	}
	sa.nextToken()

	if err := sa.parseE(); err != nil {
		return err
	}

	if !sa.initialized[varName] {
		return fmt.Errorf("переменная '%s' должна быть проинициализирована до использования", varName)
	}
	return nil
}

func (sa *SyntacticAnalyzer) parseE() error {
	if err := sa.parseT(); err != nil {
		return err
	}

	for {
		token := sa.currentToken()
		if token.Type == "Plus" {
			sa.nextToken()
			if err := sa.parseT(); err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (sa *SyntacticAnalyzer) parseT() error {
	if err := sa.parseF(); err != nil {
		return err
	}

	for {
		token := sa.currentToken()
		if token.Type == "Multiply" {
			sa.nextToken()
			if err := sa.parseF(); err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (sa *SyntacticAnalyzer) parseF() error {
	token := sa.currentToken()

	if token.Type == "Number" {
		sa.nextToken() // Число корректное, переходим к следующему токену
		return nil
	}

	if token.Type == "Identifier" {
		if !sa.initialized[token.Value] {
			return fmt.Errorf("переменная '%s' не проинициализирована", token.Value)
		}
		sa.nextToken()
		return nil
	}

	if token.Type == "LeftParen" {
		sa.nextToken()
		if err := sa.parseE(); err != nil {
			return err
		}

		token = sa.currentToken()
		if token.Type != "RightParen" {
			return fmt.Errorf("ожидалась ')', но найдено '%s'", token.Value)
		}
		sa.nextToken()
		return nil
	}

	return fmt.Errorf("ожидалось число, идентификатор или '(', но найдено '%s'", token.Value)
}
