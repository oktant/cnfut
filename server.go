package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/necais/cnfut/entities"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/", func(c echo.Context) error {
		sd := new(entities.SourceDestination)
		fmt.Println(sd)
		if err := c.Bind(sd); err != nil {
			fmt.Println(sd)
			return err
		}
		fmt.Println("dddd")
		if err := c.Validate(sd); err != nil {
			fmt.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, true)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
