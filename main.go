package main

import (
	"gql-lexer-cache/graphql"
	"gql-lexer-cache/resolver"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {

	gqlHandler := graphql.GraphQLHandler{}

	resolver.ServiceConnection = resolver.New()
	e := echo.New()
	// gql := e.Group("/graphql")
	e.POST("/gql", gqlHandler.Query)
	e.GET("/gql", gqlHandler.Playground)
	logrus.Fatal(e.Start(":8080"))
}
