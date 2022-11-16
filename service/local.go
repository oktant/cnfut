package service

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/necais/cnfut/entities"
	"github.com/necais/cnfut/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"strings"
)

import (
	"context"
	cp "github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
)

func FromLocalToS3(srcDest *entities.SourceDestination) error {
	awsSession, err := utils.GetS3Client(srcDest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	file, err := os.Open(srcDest.Source)
	if err != nil {
		return err
	}
	defer file.Close()

	// get the file size and read
	// the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	e, s3err := s3.New(awsSession).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(srcDest.Bucket),
		Key:                  aws.String(srcDest.Source),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if s3err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Info().Msg(e.String())
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
	log.Info().Msg("Successfully copied objects")
	return nil
}

func copySrcFileToDestFolder(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	defer sourceFile.Close()

	_, file := filepath.Split(src)
	newFile, err := os.Create(dest + file)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Str("Copied %d bytes.", strconv.FormatInt(bytesCopied, 10))
	return nil
}

func FromLocalToGoogle(srcDest *entities.SourceDestination) error {
	if len(strings.TrimSpace(srcDest.Bucket)) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Bucket name can't be empty")
	}
	var destObject string
	ctx := context.Background()
	client, err := utils.GetGoogleClient(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer client.Close()
	isSourceDirectory, err := utils.IsDirectory(srcDest.Source)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if isSourceDirectory {
		err := filepath.Walk(srcDest.Source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error().Msg(err.Error())
			}
			if !info.IsDir() {
				destObject := strings.Replace(path, srcDest.Source, srcDest.Destination, -1)
				err := readSourceFileAndConvertToBuffer(path, ctx, srcDest.Bucket, destObject)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			}
			return nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
