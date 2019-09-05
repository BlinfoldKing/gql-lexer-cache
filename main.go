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

	gql := e.Group("/graphql")
	// gql.Use(gqlHandler.ReadQuery)
	gql.POST("", gqlHandler.Query)
	// e.GET("/graphql", gqlHandler.Playground)
	logrus.Fatal(e.Start(":8080"))
}
