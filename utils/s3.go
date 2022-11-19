package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
	"github.com/necais/cnfut/entities"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	"strings"
)

const DefaultRegion = "us-east-1"

func GetS3Client(srcDest *entities.SourceDestination) (*session.Session, error) {
	if !validateS3Credentials(srcDest.S3AccessKeyId, srcDest.S3SecretAccessKey) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Access and Secret should be provided for s3")
	}
	region := retrieveRegion(srcDest.Region)
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			srcDest.S3AccessKeyId,
			srcDest.S3SecretAccessKey,
			""),
		Endpoint: aws.String(srcDest.Endpoint),
	})

	if err != nil {
		log.Error().Msg(err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return s, err
}

func validateS3Credentials(accessKey, secretKey string) bool {
	if len(strings.TrimSpace(accessKey)) > 1 && len(strings.TrimSpace(secretKey)) > 1 {
		return true
	}
	return false
}

func retrieveRegion(region string) string {
	if len(strings.TrimSpace(region)) > 1 {
		return region
	} else if len(strings.TrimSpace(os.Getenv("S3_REGION"))) > 1 {
		return os.Getenv("AWS_REGION")
	} else {
		return DefaultRegion
	}

}
