package filehandler

import (
	"os"
	"path/filepath"
)

// Find and process all files in a given path
func ReadFiles(path string) ([]string, error) {
	return processPath(path)
}

// Recursive function to find files
func processPath(path string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	paths := []string{}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		// If the path is a directory, list its contents and process them
		fileInfos, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, file := range fileInfos {
			subPath := filepath.Join(path, file.Name())
			subPaths, err := processPath(subPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		}

	case mode.IsRegular():
		// If the path is a regular file, validate it
		validPath := validatePath(path)
		if validPath != "" {
			paths = append(paths, validPath)
		}
	}

	return paths, nil
}

// Validate that the file is of correct file type
func validatePath(path string) string {
	extension := filepath.Ext(path)
	if extension == ".md" || extension == ".pdf" {
		return path
	}
	return ""
}

// GetNotes gets all the content from a file
func GetNotes(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), err
}
