package run

import (
	lexicalanalyzer "gitgub.com/aswait/go-syntactic-analyzer/pkg/lexical-analyzer"
	sourcetext "gitgub.com/aswait/go-syntactic-analyzer/pkg/source-text"
	screenform "gitgub.com/aswait/go-syntactic-analyzer/ui"
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
