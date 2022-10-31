package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/necais/cnfut/entities"
	"github.com/necais/cnfut/service"
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
			return err
		}
		fmt.Println("dddd")
		if err := c.Validate(sd); err != nil {
			fmt.Println(err)
			return err
		}
		copy(sd)
		return c.JSON(http.StatusOK, true)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func copy(srcDest *entities.SourceDestination) {
	switch srcDest.SourceType {
	case "s3":
		if srcDest.DestinationType == "azure" {
			service.FromS3ToAzure(srcDest)
		} else if srcDest.DestinationType == "local" {
			service.FromS3ToLocal(srcDest)
		} else if srcDest.DestinationType == "google" {
			service.FromS3ToGoogle(srcDest)
		} else {
			service.FromS3ToS3(srcDest)
		}
	case "local":
		if srcDest.DestinationType == "azure" {
			service.FromLocalToAzure(srcDest)
		} else if srcDest.DestinationType == "local" {
			service.FromLocalToLocal(srcDest)
		} else if srcDest.DestinationType == "google" {
			service.FromLocalToS3(srcDest)
		} else {
			service.FromLocalToGoogle(srcDest)
		}
	case "azure":
		if srcDest.DestinationType == "azure" {
			service.FromAzureToAzure(srcDest)
		} else if srcDest.DestinationType == "local" {
			service.FromAzureToLocal(srcDest)
		} else if srcDest.DestinationType == "google" {
			service.FromAzureToGoogle(srcDest)
		} else {
			service.FromAzureToS3(srcDest)
		}
	case "google":
		if srcDest.DestinationType == "azure" {
			service.FromGoogleToAzure(srcDest)
		} else if srcDest.DestinationType == "local" {
			service.FromGoogleToLocal(srcDest)
		} else if srcDest.DestinationType == "google" {
			service.FromGoogleToGoogle(srcDest)
		} else {
			service.FromGoogleToS3(srcDest)
		}
	}
}

