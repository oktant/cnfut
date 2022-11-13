package utils

import (
	"errors"
	zlog "github.com/rs/zerolog/log"
	"os"
)

const (
	bucketName = "your-bucket-name" // FILL IN WITH YOURS
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func IsSourceAndDestinationFolders(src, dest string) (int, error) {
	srcIsDir, err := IsDirectory(src)
	if err != nil {
		zlog.Error().Msg(err.Error())
		return 0, err
	}
	destIsDir, err := IsDirectory(dest)
	if err != nil {
		zlog.Error().Msg(err.Error())
		_, err := os.Create(dest)
		if err != nil {
			zlog.Error().Msg(err.Error())
			return 0, err
		}
	}
	if srcIsDir {
		if destIsDir {
			zlog.Info().Msg("Both src and dest are folders")
			return 1, nil
		} else {
			zlog.Error().Msg("Src is a folder and dest is a file")
			return 0, errors.New("src is a folder and dest is a file")
		}
	} else {
		if destIsDir {
			zlog.Info().Msg("Src is a file and dest is a folder")
			return 3, nil
		} else {
			zlog.Info().Msg("Src is a file and dest is a file")
			return 4, nil
		}

	}
}
