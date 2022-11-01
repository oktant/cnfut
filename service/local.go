package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/necais/cnfut/utils"
	"io"
	"net/http"
	"strconv"

	//"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/necais/cnfut/entities"
	zlog "github.com/rs/zerolog/log"
	"os"
)

func FromLocalToS3(srcDest *entities.SourceDestination) {
	//s3Client, err := minio.New("s3.amazonaws.com", &minio.Options{
	//	Creds:  credentials.NewStaticV4("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
	//	Secure: true,
	//})

}

func FromLocalToAzure(srcDest *entities.SourceDestination) {}

func FromLocalToLocal(srcDest *entities.SourceDestination) error {
	folderOrFile, err := utils.IsSourceAndDestinationFolders(srcDest.Source, srcDest.Destination)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if folderOrFile == 1 {
		copySrcFolderToDestFolder(srcDest.Source, srcDest.Destination)
	} else if folderOrFile == 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Destination is a file while source is a folder")
	} else if folderOrFile == 3 {
		copySrcFileToDestFolder(srcDest.Source, srcDest.Destination)
	} else {
		err := copySrcFileToDestFile(srcDest.Source, srcDest.Destination)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return nil
}

func copySrcFolderToDestFolder(src, dest string) {

}

func copySrcFileToDestFolder(src, dest string) {

}

func copySrcFileToDestFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		zlog.Error().Err(err)

		return err
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create("/var/www/html/test.txt")
	if err != nil {
		zlog.Error().Err(err)
		return err
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		zlog.Error().Err(err)
		return err
	}
	zlog.Info().Str("Copied %d bytes.", strconv.FormatInt(bytesCopied, 10))
	return nil
}

func FromLocalToGoogle(srcDest *entities.SourceDestination) {}

func getAWSSession(region string) *session.Session {
	awsConfig := &aws.Config{}
	if len(region) > 0 {
		awsConfig.Region = aws.String(region)
	} else if len(os.Getenv("AWS_REGION")) > 0 {
		awsConfig.Region = aws.String(region)
	} else {
		awsConfig.Region = aws.String(entities.DefaultAwsRegion)
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region)},
	)
	if err != nil {
		panic(err)
	}
	return sess
}
