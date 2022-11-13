package utils

import (
	"errors"
	"github.com/necais/cnfut/entities"
	zlog "github.com/rs/zerolog/log"
	"os"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func IsSourceAndDestinationFolders(srcDest *entities.SourceDestination) (int, error) {
	srcIsDir, err := IsDirectory(srcDest.Source)
	if err != nil {
		zlog.Error().Msg(err.Error())
		return 0, err
	}
	destIsDir, err := IsDirectory(srcDest.Destination)
	if err != nil {
		zlog.Error().Msg(err.Error())
		_, err := os.Create(srcDest.Destination)
		if err != nil {
			zlog.Error().Msg(err.Error())
			return 0, err
		}
	}
	if srcIsDir {
		if destIsDir {
			AddPathSeparatorToFolders(srcDest.Source)
			AddPathSeparatorToFolders(srcDest.Destination)
			zlog.Info().Msg("Both src and dest are folders")
			return 1, nil
		} else {
			AddPathSeparatorToFolders(srcDest.Source)
			zlog.Error().Msg("Src is a folder and dest is a file")
			return 0, errors.New("src is a folder and dest is a file")
		}
	} else {
		if destIsDir {
			AddPathSeparatorToFolders(srcDest.Destination)
			zlog.Info().Msg("Src is a file and dest is a folder")
			return 3, nil
		} else {
			zlog.Info().Msg("Src is a file and dest is a file")
			return 4, nil
		}

	}
}

func AddPathSeparatorToFolders(folder string) string {
	if folder[len(folder)-1:] != string(os.PathSeparator) {
		folder = folder + string(os.PathSeparator)
	}
	return folder
}
