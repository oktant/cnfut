package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/necais/cnfut/entities"
	"github.com/rs/zerolog/log"
	"os"

	"strings"
)

const DefaultRegion = "eu-north-1"

func GetS3Client(srcDest *entities.SourceDestination) (*session.Session, error) {
	region := retrieveRegion(srcDest.Region)
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			"Q3AM3UQ867SPQQA43P2F",
			"zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
			""),
		Endpoint: aws.String("play.min.io"),
	})

	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return s, err
}

func retrieveRegion(region string) string {
	if len(strings.TrimSpace(region)) > 1 {
		return region
	} else if len(strings.TrimSpace(os.Getenv("AWS_REGION"))) > 1 {
		return os.Getenv("AWS_REGION")
	} else {
		return DefaultRegion
	}

}
