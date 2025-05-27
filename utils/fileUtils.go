package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAbsolutePath(file string) string {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return file
	}
	return absPath
}

func WriteToFile(file string, content string, flag int) error {
	outFile, err := os.OpenFile(file, flag, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer outFile.Close()

	_, err = outFile.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

func DeleteFileIfExists(file string) {
	if _, err := os.Stat(file); err == nil {
		err = os.Remove(file)
		if err != nil {
			fmt.Printf("Error deleting file %s: %v\n", file, err)
		}
	}
}
