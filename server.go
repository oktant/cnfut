package main

import (
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
		if err := c.Bind(sd); err != nil {
			return err
		}
		if err := c.Validate(sd); err != nil {
			return err
		}
		err := copyObject(sd)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, true)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func copyObject(srcDest *entities.SourceDestination) error {
	switch srcDest.SourceType {
	case "s3":
		if srcDest.DestinationType == "azure" {
			err := service.FromS3ToAzure(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "local" {
			err := service.FromS3ToLocal(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "google" {
			err := service.FromS3ToGoogle(srcDest)
			if err != nil {
				return err
			}
		} else {
			err := service.FromS3ToS3(srcDest)
			if err != nil {
				return err
			}
		}
	case "local":
		if srcDest.DestinationType == "azure" {
			err := service.FromLocalToAzure(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "local" {
			err := service.FromLocalToLocal(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "google" {
			err := service.FromLocalToS3(srcDest)
			if err != nil {
				return err
			}
		} else {
			service.FromLocalToGoogle(srcDest)
		}
	case "azure":
		if srcDest.DestinationType == "azure" {
			err := service.FromAzureToAzure(srcDest)
			if err != nil {
				return err
			}

		} else if srcDest.DestinationType == "local" {
			err := service.FromAzureToLocal(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "google" {
			err := service.FromAzureToGoogle(srcDest)
			if err != nil {
				return err
			}
		} else {
			err := service.FromAzureToS3(srcDest)
			if err != nil {
				return err
			}
		}
	case "google":
		if srcDest.DestinationType == "azure" {
			err := service.FromGoogleToAzure(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "local" {
			err := service.FromGoogleToLocal(srcDest)
			if err != nil {
				return err
			}
		} else if srcDest.DestinationType == "google" {
			err := service.FromGoogleToGoogle(srcDest)
			if err != nil {
				return err
			}
		} else {
			err := service.FromGoogleToS3(srcDest)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
