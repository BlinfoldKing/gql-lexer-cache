package lexer

import (
	"fmt"
	"gql-lexer-cache/graphql/lexer/token/operation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQueryOperation(t *testing.T) {
	t.Run("should return correct operation name and type of named operation", func(t *testing.T) {

		query := `query IniOperation{
				helloWorld(), halodunia, recursiveQuery(){
					props1
				}
				another
		}`

		lexer := New(query)
		op, err := lexer.Parse()
		fmt.Println(op, err)

		assert.NoError(t, err)
		assert.Equal(t, operation.Query, op[0].Type)
		assert.Equal(t, "IniOperation", op[0].Name)
		assert.Equal(t, "helloWorld", op[0].Selections[0].Name)
		assert.Equal(t, "halodunia", op[0].Selections[1].Name)
		assert.Equal(t, "recursiveQuery", op[0].Selections[2].Name)
		subSelection := op[0].Selections[2].Value
		assert.Equal(t, "props1", subSelection[0].Name)
	})

	t.Run("should able to get anonymous operation", func(t *testing.T) {
		query := `query {
			helloWorld
		}`
		lexer := New(query)
		op, err := lexer.Parse()
		fmt.Println(op)
		assert.Nil(t, err)
		assert.Equal(t, operation.Query, op[0].Type)
		assert.Equal(t, "", op[0].Name)
	})

	t.Run("should able to get anonymous operation", func(t *testing.T) {
		query := `{
			helloWorld
		}`
		lexer := New(query)
		op, err := lexer.Parse()
		fmt.Println(op)
		assert.Nil(t, err)
		assert.Equal(t, operation.Query, op[0].Type)
		assert.Equal(t, "", op[0].Name)
	})
}

func TestParseMutationOperation(t *testing.T) {
	t.Run("should return correct operation name and type of named operation", func(t *testing.T) {
		query := `mutation IniOperation{
				helloWorld(), halodunia, recursiveQuery(){
					props1
				}
				another
		}`

		lexer := New(query)
		op, err := lexer.Parse()
		fmt.Println(op, err)

		assert.NoError(t, err)
		assert.Equal(t, operation.Mutation, op[0].Type)
		assert.Equal(t, "IniOperation", op[0].Name)
		assert.Equal(t, "helloWorld", op[0].Selections[0].Name)
		assert.Equal(t, "halodunia", op[0].Selections[1].Name)
		assert.Equal(t, "recursiveQuery", op[0].Selections[2].Name)
		subSelection := op[0].Selections[2].Value
		assert.Equal(t, "props1", subSelection[0].Name)
	})

	t.Run("should able to get anonymous operation", func(t *testing.T) {
		query := `mutation {
			helloWorld
		}`
		lexer := New(query)
		op, err := lexer.Parse()
		fmt.Println(op)
		assert.Nil(t, err)
		assert.Equal(t, operation.Mutation, op[0].Type)
		assert.Equal(t, "", op[0].Name)
	})

}
