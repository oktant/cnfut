package service

import (
	"bytes"
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
	var destObject string

	isSourceDirectory, err := utils.IsDirectory(srcDest.Source)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	awsSession, err := utils.GetS3Client(srcDest)
	if err != nil {
		return err
	}
	if isSourceDirectory {
		err := filepath.Walk(srcDest.Source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error().Msg(err.Error())
			}
			if !info.IsDir() {
				destObject := strings.Replace(path, srcDest.Source, srcDest.Destination, -1)
				err := utils.UploadAFileToS3(path, awsSession, srcDest.Bucket, destObject)
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
		err := utils.UploadAFileToS3(srcDest.Source, awsSession, srcDest.Bucket, destObject)
		if err != nil {
			return err
		}

	}

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
