package sourcetext

import "errors"

type SourceTexter interface {
	ReadNextSymbol() (rune, error)
	ResetCursor()
	HasMoreSymbols() bool
	LoadFromInput(input string)
}

type SourceText struct {
	text   []rune
	cursor int
}

func NewSourceText() *SourceText {
	return &SourceText{
		text:   []rune{},
		cursor: 0,
	}
}

func (st *SourceText) ReadNextSymbol() (rune, error) {
	if st.cursor >= len(st.text) {
		return 0, errors.New("Gопытка чтения за пределами исходного текста")
	}

	symbol := rune(st.text[st.cursor])
	st.cursor++
	return symbol, nil
}

func (st *SourceText) ResetCursor() {
	st.cursor = 0
}

func (st *SourceText) HasMoreSymbols() bool {
	return st.cursor < len(st.text)
}

func (st *SourceText) LoadFromInput(input string) {
	st.text = []rune(input)
	st.cursor = 0
}
