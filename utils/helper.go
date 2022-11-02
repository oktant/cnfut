package utils

import (
	"fmt"
	"os"
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
		return 0, err
	}
	destIsDir, err := IsDirectory(dest)
	if err != nil {
		_, err := os.Create(dest)
		if err != nil {
			return 0, err
		}
	}
	if srcIsDir {
		if destIsDir {
			fmt.Println("Both src and dest are folders")
			return 1, nil
		} else {
			return 2, nil
		}
	} else {
		if destIsDir {
			fmt.Println("Src is a file and dest is a folder")
			return 3, nil
		} else {
			fmt.Println("Src is a file and dest is a file")
			return 4, nil
		}

	}
}
