package gitsupports

import (
	"fmt"
	"merge/constants"
	"merge/objects"
	"os/exec"
	"strings"
)

func GetChangedFiles(commitHash string, config objects.Config) {
	// Run the git command to get the list of changed files
	cmd := exec.Command("git", "show", "--pretty=", "--name-only", commitHash)
	cmd.Dir = config.GitRepo
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// Split the output into lines and return as a slice of strings
	files := strings.Split(string(output), constants.BreakLine)
	fmt.Println(constants.Empty)
	// concat files with git repo to get full output path
	for i, file := range files {
		if file != constants.Empty {
			files[i] = config.GitRepo + constants.PathSeparator + file
			// fmt.Println(files[i])
		}
	}

	orderedFiles := orderedFiles(files)

	for _, file := range orderedFiles {
		fmt.Println(file)
	}

	fmt.Println(constants.Empty)
}

func orderedFiles(files []string) []string {
	// Sort the files in alphabetical order
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i] > files[j] {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
	return files
}

func ExecGitCommand(command string, config objects.Config) {
	// Run the git command to get the list of changed files
	cmd := exec.Command("git", command)
	cmd.Dir = config.GitRepo
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// Print the output
	fmt.Println(string(output))

}
