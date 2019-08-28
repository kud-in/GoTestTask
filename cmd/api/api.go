package main

import (
	"TestTask/handlers"
	"TestTask/models"
	"database/sql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	e := echo.New()
	e.Logger.SetLevel(log.OFF)

	e.Use(middleware.Logger())  // logger middleware will “wrap” recovery
	e.Use(middleware.Recover()) // as it is enumerated before in the Use calls

	e.File("/", "static/swagger.html")
	e.File("/validator", "static/validator.png")

	e.Static("/static", "static")

	v1 := e.Group("/v1")

	v1.GET("/status/:code", handlers.GetStatus)

	v1.GET("/status", handlers.GetStatus)

	v1.GET("/cost/:code", handlers.GetCost)

	v1.GET("/history/:code", handlers.GetHistory)

	db, err := sql.Open("sqlite3", "go_test.sqlite")
	if err != nil {
		panic(err.Error())

	} else {
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set(models.DBContextKey, db)

				return next(c)
			}
		})
	}

	e.Logger.Fatal(e.Start(":8080"))
}
