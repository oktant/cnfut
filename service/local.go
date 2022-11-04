package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
	"github.com/necais/cnfut/entities"
	"github.com/necais/cnfut/utils"
	cp "github.com/otiai10/copy"
	zlog "github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func FromLocalToS3(srcDest *entities.SourceDestination) error {
	//s3Client, err := minio.New("s3.amazonaws.com", &minio.Options{
	//	Creds:  credentials.NewStaticV4("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
	//	Secure: true,
	//})
	return nil
}

func FromLocalToAzure(srcDest *entities.SourceDestination) error {
	return nil
}

func FromLocalToLocal(srcDest *entities.SourceDestination) error {
	folderOrFile, err := utils.IsSourceAndDestinationFolders(srcDest.Source, srcDest.Destination)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if folderOrFile == 1 || folderOrFile == 4 {
		err := cp.Copy(srcDest.Source, srcDest.Destination)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else if folderOrFile == 3 {
		err := copySrcFileToDestFolder(srcDest.Source, srcDest.Destination)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	zlog.Info().Msg("Successfully copied objects")
	return nil
}

func copySrcFileToDestFolder(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		zlog.Error().Err(err)
		return err
	}
	defer sourceFile.Close()

	_, file := filepath.Split(src)
	if dest[len(dest)-1:] != string(os.PathSeparator) {
		dest = dest + string(os.PathSeparator)
	}

	newFile, err := os.Create(dest + file)
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
