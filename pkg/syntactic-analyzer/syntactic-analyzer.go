package syntacticanalyzer

type SyntacticAnalyzerer interface {
}

type SyntacticAnalyzer struct {
}

func NewSyntacticAnalyzer() *SyntacticAnalyzer {
	return &SyntacticAnalyzer{}
}
