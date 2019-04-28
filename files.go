package goutils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/jatgam/goutils/log"
)

// DirSizeMB returns the size of the provided directory in MB. Doesn't traverse
// into subdirectories.
func DirSizeMB(filePath string) (float64, error) {
	dirSizeBytes, err := DirSizeBytes(filePath)
	if err != nil {
		return 0, err
	}
	sizeMB := float64(dirSizeBytes) / 1024.0 / 1024.0

	return sizeMB, nil
}

// DirSizeBytes returns the size of the provided directory in MB. Doesn't traverse
// into subdirectories.
func DirSizeBytes(filePath string) (int64, error) {
	var dirSize int64

	readSize := func(filePath string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return err
	}

	err := filepath.Walk(filePath, readSize)
	return dirSize, err
}

// DeleteOldestFileUntilDirUnderSize will delete the oldest files in the supplied
// path until the directory size is under the supplied size in MB
func DeleteOldestFileUntilDirUnderSize(dirPath string, sizeLimit float64) error {
	var allFoundFiles []os.FileInfo
	var filesToDelete []os.FileInfo
	pathStat, statErr := os.Stat(dirPath)
	if statErr != nil {
		return statErr
	}
	if !pathStat.IsDir() {
		return fmt.Errorf("Supplied path isn't a directory: %s", dirPath)
	}
	err := filepath.Walk(dirPath, func(dirPath string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			allFoundFiles = append(allFoundFiles, file)
		}
		return err
	})
	if err != nil {
		return err
	}
	sort.Slice(allFoundFiles, func(i int, j int) bool {
		return allFoundFiles[i].ModTime().Unix() < allFoundFiles[j].ModTime().Unix() //oldest first
	})
	currentDirSize, err := DirSizeMB(dirPath)
	if err != nil {
		return fmt.Errorf("Couldn't determine directories current size: %s", dirPath)
	}
	if currentDirSize < sizeLimit {
		return nil
	}
	for _, fileFound := range allFoundFiles {
		if currentDirSize < sizeLimit {
			break
		}
		fileSizeMB := float64(fileFound.Size()) / 1024.0 / 1024.0
		currentDirSize -= fileSizeMB
		filesToDelete = append(filesToDelete, fileFound)
	}
	for _, fileTD := range filesToDelete {
		removeErr := os.Remove(filepath.FromSlash(path.Join(dirPath, fileTD.Name())))
		if removeErr == nil {
			log.WithFields(log.Fields{
				"File": fileTD.Name(),
				"Size": fmt.Sprintf("%.2f", float64(fileTD.Size())/1024.0/1024.0),
			}).Debug("Deleted file from Directory")
		} else {
			log.WithFields(log.Fields{
				"File": fileTD.Name(),
				"Size": fmt.Sprintf("%.2f", float64(fileTD.Size())/1024.0/1024.0),
			}).Error("Failed to delete file from cache")
		}

	}
	return nil
}

// UpdateFileModifiedTime will set the atime and mtime to the supplied time on the
// supplied file
func UpdateFileModifiedTime(filePath string, timeToSet time.Time) error {
	_, err := os.Stat(filePath)
	if err == nil {
		currenttime := timeToSet.Local()
		updateErr := os.Chtimes(filePath, currenttime, currenttime)
		if updateErr != nil {
			return updateErr
		}
	} else {
		return fmt.Errorf("File doesn't exist: %s", filePath)
	}
	return nil
}
