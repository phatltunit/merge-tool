package utils

import "path/filepath"

func GetAbsolutePath(file string) string {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return file
	}
	return absPath
}
