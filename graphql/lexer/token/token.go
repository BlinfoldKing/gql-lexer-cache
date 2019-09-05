package token

import (
	"gql-lexer-cache/graphql/lexer/token/operation"
	"gql-lexer-cache/graphql/lexer/token/selection"
)

type (
	SelectionSet []Selection
)

// Operation = selection_set | operation_type NAME variable_definitions? directives? selection_set
type Operation struct {
	Type operation.Type
	Name string
	// TODO: not yet implemented
	// Variables  []string
	// Directives []string
	Selections SelectionSet
}

// containing selection information that
// could be a single field selection
// or recursive selection field
type Selection struct {
	// "field" | "selection",
	// use to refelct the
	// type of selection value
	Type selection.Type
	// nil | SelectionSet
	Value SelectionSet
	// name of the query
	Name string
}
