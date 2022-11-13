package service

import (
	"bytes"
	"context"
	"fmt"
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
	"time"
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

func FromLocalToGoogle(srcDest *entities.SourceDestination) error {
	var destObject string
	ctx := context.Background()
	client, err := utils.GetInstance(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	folderOrFile, err := utils.IsSourceAndDestinationFolders(srcDest.Source, srcDest.Destination)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if folderOrFile == 4 {
		destObject = filepath.Dir(srcDest.Destination) + string(os.PathSeparator) + filepath.Base(srcDest.Source)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else if folderOrFile == 1 {
		err := copySrcFileToDestFolder(srcDest.Source, srcDest.Destination)
		if err != nil {
			///ToDo add folder support
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else if folderOrFile == 3 {
		if srcDest.Destination[len(srcDest.Destination)-1:] != string(os.PathSeparator) {
			srcDest.Destination = srcDest.Destination + string(os.PathSeparator)
		}
		destObject = srcDest.Destination + filepath.Base(srcDest.Source)
	}

	fileContent, err := os.ReadFile(srcDest.Source)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	buf := bytes.NewBuffer(fileContent)
	err = uploadFileToGoogle(ctx, buf, srcDest.Bucket, destObject)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil

}

func uploadFileToGoogle(ctx context.Context, buf *bytes.Buffer, bucket, object string) error {
	client, err := utils.GetInstance(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 5

	if _, err = io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	zlog.Info().Msg("Successfully copied objects")
	return nil

}

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
