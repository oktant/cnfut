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
	"strings"
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
	folderOrFile, err := utils.IsSourceAndDestinationFolders(srcDest)
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
	isSourceDirectory, err := utils.IsDirectory(srcDest.Source)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if isSourceDirectory {
		err := filepath.Walk(srcDest.Source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				zlog.Error().Msg(err.Error())
			}
			if !info.IsDir() {
				destObject := strings.Replace(path, srcDest.Source, srcDest.Destination, -1)
				err := readSourceFileAndConvertToBuffer(path, ctx, srcDest.Bucket, destObject)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		destObject = srcDest.Destination + string(os.PathSeparator) + filepath.Base(srcDest.Source)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		err := readSourceFileAndConvertToBuffer(srcDest.Source, ctx, srcDest.Bucket, destObject)
		if err != nil {
			return err
		}

	}

	return nil

}

//func copyFolderToGoogleCloud(source, destination string) error {
//	files, err := os.ReadDir(source)
//	if err != nil {
//		zlog.Error().Msg(err.Error())
//		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//	}
//	for _, file := range files {
//		fmt.Println(file.Name(), file.IsDir())
//	}
//}

func readSourceFileAndConvertToBuffer(source string, ctx context.Context, bucket, destObject string) error {
	fileContent, err := os.ReadFile(source)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	buf := bytes.NewBuffer(fileContent)
	err = utils.UploadFileToGoogle(ctx, buf, bucket, destObject)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
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
