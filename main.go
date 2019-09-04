package main

import (
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()
	e.GET("/", func(ctx echo.Context) error {
		ctx.JSON(200, echo.Map{
			"hello": "world",
		})
		return nil
	})
	logrus.Fatal(e.Start(":8080"))
}
