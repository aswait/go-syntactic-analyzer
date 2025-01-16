package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	lexicalanalyzer "gitgub.com/aswait/go-syntactic-analyzer/pkg/lexical-analyzer"
	syntacticanalyzer "gitgub.com/aswait/go-syntactic-analyzer/pkg/syntactic-analyzer"
)

type ScreenFormer interface {
	Run()
}

type ScreenForm struct {
	App             fyne.App
	Window          fyne.Window
	InputField      *widget.Entry
	OutputField     *widget.Label
	StartButton     *widget.Button
	LexicalAnalyzer lexicalanalyzer.LexicalAnalyzerer
	SyntaxAnalyzer  syntacticanalyzer.SyntacticAnalyzer
}

func NewScreenForm(lexicalAnalyzer lexicalanalyzer.LexicalAnalyzerer, syntaxAnalyzer syntacticanalyzer.SyntacticAnalyzer) *ScreenForm {
	myApp := app.New()
	myWindow := myApp.NewWindow("Лексический анализатор")

	myApp.Settings().SetTheme(theme.LightTheme())

	inputField := widget.NewMultiLineEntry()
	inputField.SetPlaceHolder("Введите текст для анализа...")

	outputField := widget.NewLabel("")
	outputField.Wrapping = fyne.TextWrapWord

	sf := &ScreenForm{
		App:             myApp,
		Window:          myWindow,
		InputField:      inputField,
		OutputField:     outputField,
		LexicalAnalyzer: lexicalAnalyzer,
		SyntaxAnalyzer:  syntaxAnalyzer,
	}

	sf.StartButton = widget.NewButton("Запуск", func() {
		inputText := sf.InputField.Text
		if inputText == "" {
			sf.OutputField.SetText("Ошибка: текст для анализа пустой.")
			return
		}

		sf.LexicalAnalyzer.SourceLoadFromInput(inputText)
		_, err := sf.LexicalAnalyzer.Transliterate()
		if err != nil {
			sf.OutputField.SetText(fmt.Sprintf("Ошибка: %v", err))
			return
		}

		tokens, err := sf.LexicalAnalyzer.Analyze()
		if err != nil {
			sf.OutputField.SetText(fmt.Sprintf("Ошибка: %v", err))
			return
		}

		fmt.Println(tokens)

		err = sf.SyntaxAnalyzer.Parse(tokens)
		if err != nil {
			sf.OutputField.SetText(fmt.Sprintf("Ошибка: %v", err))
			return
		}

		finalMessage := sf.LexicalAnalyzer.Validate()

		sf.OutputField.SetText("\n" + finalMessage)
	})

	content := container.NewVBox(
		widget.NewLabel("Введите текст:"),
		sf.InputField,
		sf.StartButton,
		widget.NewLabel("Результаты:"),
		sf.OutputField,
	)

	sf.Window.SetContent(content)
	sf.Window.Resize(fyne.NewSize(600, 400))
	return sf
}

func (sf *ScreenForm) Run() {
	sf.Window.ShowAndRun()
}
