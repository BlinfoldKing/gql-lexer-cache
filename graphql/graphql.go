package graphql

import (
	"encoding/json"
	"gql-lexer-cache/graphql/lexer"
	"gql-lexer-cache/resolver"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type GraphQLHandler struct {
	Lexer lexer.Lexer
}

// GET server graphql playground
func (handler *GraphQLHandler) Playground(c echo.Context) error {
	c.HTML(http.StatusOK, page)
	return nil
}

// POST graphql query
func (handler *GraphQLHandler) Query(c echo.Context) error {
	s, err := getSchema("schema/schema.gql")
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = GqlResponse(c, graphql.MustParseSchema(s, resolver.ServiceConnection))
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func GqlResponse(ctx echo.Context, Schema *graphql.Schema) error {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := json.NewDecoder(ctx.Request().Body).Decode(&params); err != nil {
		return err
	}

	response := Schema.Exec(ctx.Request().Context(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	ctx.Response().Write(responseJSON)

	return nil
}

func getSchema(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
