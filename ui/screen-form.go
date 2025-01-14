package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	lexicalanalyzer "gitgub.com/aswait/go-syntactic-analyzer/pkg/lexical-analyzer"
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
}

func NewScreenForm(lexicalanalyzer lexicalanalyzer.LexicalAnalyzerer) *ScreenForm {
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
		LexicalAnalyzer: lexicalanalyzer,
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

		result, err := sf.LexicalAnalyzer.Analyze()
		if err != nil {
			sf.OutputField.SetText(fmt.Sprintf("Ошибка: %v", err))
			return
		}

		finalMessage := sf.LexicalAnalyzer.Validate()

		sf.OutputField.SetText(strings.Join(result, "\n") + "\n" + finalMessage)
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
