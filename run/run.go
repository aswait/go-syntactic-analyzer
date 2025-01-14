package run

import (
	lexicalanalyzer "gitgub.com/aswait/go-lexical-analyzer/pkg/lexical-analyzer"
	sourcetext "gitgub.com/aswait/go-lexical-analyzer/pkg/source-text"
	screenform "gitgub.com/aswait/go-lexical-analyzer/screen-form"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	sourcetext := sourcetext.NewSourceText()

	lexicalanalyzer := lexicalanalyzer.NewLexicalAnalyzer(sourcetext)

	screenform := screenform.NewScreenForm(lexicalanalyzer)
	screenform.Run()
}
