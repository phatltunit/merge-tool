package gitsupports

import (
	"os/exec"
	"fmt"
	"strings"
	"merge/objects"
	"merge/constants"
)


func GetChangedFiles(commitHash string, config objects.Config){
	// Run the git command to get the list of changed files
	cmd := exec.Command("git", "show", "--pretty=", "--name-only", commitHash)
	cmd.Dir = config.GitRepo
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// Split the output into lines and return as a slice of strings
	files := strings.Split(string(output), constants.BreakLine)
	fmt.Println(constants.BreakLine)
	// concat files with git repo to get full output path
	for i, file := range files {
		if file != constants.Empty {
			files[i] = config.GitRepo + constants.PathSeparator + file
			fmt.Println(files[i])
		}
	}
	fmt.Println(constants.BreakLine)
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