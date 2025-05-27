package gitsupports

import (
	"fmt"
	"merge/constants"
	"merge/objects"
	"merge/utils"
	"os/exec"
	"sort"
	"strings"
	"sync"
)

// we will use goroutines to check files in parallel
func checkFileChanged(commitHash string, config objects.Config) []string {
	hash := strings.Split(commitHash, ",")
	// Filter out empty strings from the hash slice
	var filteredHash []string
	for _, h := range hash {
		if h != constants.Empty {
			filteredHash = append(filteredHash, h)
		}
	}
	hash = filteredHash

	var wg sync.WaitGroup
	fileChan := make(chan []string, len(hash))
	wg.Add(len(hash)) // add the number of goroutines to wait for
	for _, commit := range hash {
		go singleCheck(commit, config, fileChan, &wg)
	}
	wg.Wait()       // wait for all goroutines to finish
	close(fileChan) // this is important to close the channel after all goroutines are done
	// Collect all files from the channel
	var allFiles []string
	for files := range fileChan {
		allFiles = append(allFiles, files...)
	}
	return allFiles
}

func singleCheck(commitHash string, config objects.Config, fileChan chan []string, wg *sync.WaitGroup) {
	defer wg.Done() // signal that this goroutine is done
	// Run the git command to check if the file has changed
	cmd := exec.Command("git", "show", "--pretty=", "--name-only", commitHash)
	cmd.Dir = config.GitRepo
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Split the output into lines and check if the file is in the list
	files := strings.Split(string(output), constants.BreakLine)
	// concat files with git repo to get full output path
	for i, file := range files {
		if file != constants.Empty {
			files[i] = utils.GetAbsolutePath(config.GitRepo + constants.PathSeparator + file)
		}
	}

	fileChan <- files
}

func GetChangedFiles(commitHash string, config objects.Config) {
	files := checkFileChanged(commitHash, config)
	orderedFiles := orderedFiles(files)
	for _, file := range orderedFiles {
		fmt.Println(file)
	}
	fmt.Println(constants.Empty)
}

func orderedFiles(files []string) []string {
	// Sort the files in alphabetical order
	distinctFiles := make(map[string]bool, len(files))
	orderedFiles := make([]string, 0, len(files))

	for _, file := range files {
		if _, ok := distinctFiles[file]; !ok {
			distinctFiles[file] = true
			orderedFiles = append(orderedFiles, file)
		}
	}

	sort.Strings(orderedFiles)
	return orderedFiles
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
