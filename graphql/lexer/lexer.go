package lexer

// Lexer definition
type Lexer struct {
	input      string
	parseStack []rune
	cursor     int
}

// New use to init lexer
func New(gql string) (l *Lexer) {
	return &Lexer{
		input: gql,
	}
}
