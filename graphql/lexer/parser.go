package lexer

import (
	"errors"
	"gql-lexer-cache/graphql/lexer/token"
	"gql-lexer-cache/graphql/lexer/token/operation"
	"gql-lexer-cache/graphql/lexer/token/selection"
)

func isWhitespace(c rune) bool {
	switch c {
	case ' ', '\n', '\t':
		return true
	default:
		return false
	}
}

func isInlineWhitespace(c rune) bool {
	switch c {
	case ' ', '\t':
		return true
	default:
		return false
	}
}

func (l *Lexer) consumeWhitespace() {
	c, err := l.read()
	for err == nil && isWhitespace(c) {
		if c == '\n' {
			l.push('\\')
		}
		l.cursor++
		c, err = l.read()
	}
}

func (l *Lexer) consumeWhitespaceUntil(stop rune) {
	c, err := l.read()
	for err == nil && isWhitespace(c) && c != stop {
		if c == '\n' {
			l.push('\\')
		}
		l.cursor++
		c, err = l.read()
	}
}

func (l *Lexer) isEOF() bool {
	return l.cursor >= len(l.input)
}

func (l *Lexer) read() (c rune, err error) {
	if l.isEOF() {
		err = errors.New("end of file")
	} else {
		c = rune(l.input[l.cursor])
	}
	return
}

// Parse is used as the main parser
func (l *Lexer) Parse() (operations []token.Operation, err error) {
	l.cursor = 0
	for !l.isEOF() {
		l.consumeWhitespace()
		op, err := l.parseOperation()
		if err != nil {
			return nil, err
		}
		operations = append(operations, op)
		l.cursor++
		l.consumeWhitespace()
	}

	return
}

func (l *Lexer) parseOperation() (op token.Operation, err error) {
	l.pushFlush()
	// parse "query" keyword
	var (
		isQuery    bool
		isMutation bool
	)

	c, err := l.read()
	if err != nil {
		return
	}
	switch c {
	case 'q':
		if err = l.parseQueryKeyword(); err != nil {
			return
		}
		isQuery = true
		l.cursor++
	case 'm':
		if err = l.parseMutationKeyword(); err != nil {
			return
		}
		isMutation = true
		l.cursor++
	case '{':
		break
	default:
		err = errors.New("unknown definition")
		return
	}

	if isQuery {
		op.Type = operation.Query
	} else if isMutation {
		op.Type = operation.Mutation
	}

	// get name of named operation
	if isQuery || isMutation {
		op.Name = l.parseName()
		l.pushString(op.Name)
	} else {
		op.Type = operation.Query
	}

	l.consumeWhitespace()

	c, err = l.read()
	if err != nil {
		return
	}

	if c == '{' {
		l.push('{')
		l.cursor++
		op.Selections, err = l.parseSelectionSet()
		if err != nil {
			return
		}
	} else {
		err = errors.New("no selection found")
		return
	}

	return
}

func (l *Lexer) parseSelection() (newSelection token.Selection, err error) {
	l.pushFlush()
	name := l.parseName()
	if name != "" {
		newSelection = token.Selection{
			Type: selection.Field,
			Name: name,
		}
	}
	l.pushString(name)
	l.consumeWhitespace()
	c, err := l.read()
	if err != nil {
		return
	}

	if c == '(' {
		l.push('(')
		l.cursor++
		l.consumeWhitespace()
		c, err = l.read()
		if err != nil {
			return
		}
		if c != ')' {
			err = errors.New("mismatch )")
			return
		}
		l.push(')')
		l.cursor++
	}

	c, err = l.read()
	if err != nil {
		return
	}

	if c == '{' {
		l.push('{')
		l.cursor++
		subSelection, err := l.parseSelectionSet()
		if err != nil {
			return token.Selection{}, err
		}
		newSelection.Value = subSelection
	}

	c, err = l.read()

	l.popFlush()
	return
}

func (l *Lexer) parseSelectionSet() (selectionSet token.SelectionSet, err error) {
	l.pushFlush()
	l.consumeWhitespace()
	c, err := l.read()
	for err == nil && c != '}' {
		if top := l.pop(); len(selectionSet) >= 1 && top != '\\' && top != ',' {
			return token.SelectionSet{}, errors.New(`no separator found`)
		}
		newSelection, err := l.parseSelection()
		if err != nil {
			return token.SelectionSet{}, err
		}
		selectionSet = append(selectionSet, newSelection)
		l.consumeWhitespaceUntil(',')
		if c, err := l.read(); err == nil && c == ',' {
			l.push(',')
			l.cursor++
		}
		l.consumeWhitespace()
		c, err = l.read()
	}
	l.cursor++
	err = l.popFlush()
	return
}

func (l *Lexer) parseQueryKeyword() error {
	var err error
	l.parseStack = append(l.parseStack, 'q')
	l.cursor++
	for !isWhitespace(rune(l.input[l.cursor])) && !l.isEOF() {
		switch l.input[l.cursor] {
		case 'u':
			err = l.popCond('q')
		case 'e':
			err = l.popCond('u')
		case 'r':
			err = l.popCond('e')
		case 'y':
			err = l.popCond('r')
		default:
			return errors.New("unknown keyword")
		}
		if err != nil {
			return err
		}
		l.push(rune(l.input[l.cursor]))
		l.cursor++
	}
	err = l.popCond('y')
	return err
}

func (l *Lexer) parseMutationKeyword() error {
	var err error
	l.parseStack = append(l.parseStack, 'm')
	l.cursor++
	for !isWhitespace(rune(l.input[l.cursor])) && !l.isEOF() {
		switch l.input[l.cursor] {
		case 'u':
			err = l.popCond('m')
		case 't':
			prev := l.pop()
			if prev != 'u' && prev != 'a' {
				err = errors.New("unknown keyword")
			}
		case 'a':
			err = l.popCond('t')
		case 'i':
			err = l.popCond('t')
		case 'o':
			err = l.popCond('i')
		case 'n':
			err = l.popCond('o')
		default:
			return errors.New("unknown keyword")
		}
		if err != nil {
			return err
		}
		l.push(rune(l.input[l.cursor]))
		l.cursor++
	}
	err = l.popCond('n')
	return err
}

func (l *Lexer) parseName() string {
	var name string
	c, err := l.read()
	for !isWhitespace(c) &&
		c != '{' &&
		c != '}' &&
		c != '(' &&
		c != ')' &&
		c != ',' &&
		err == nil {
		name += string(c)
		l.cursor++
		c, err = l.read()
	}

	return name
}
