package logics

import (
	"bufio"
	"flag"
	"fmt"
	"merge/constants"
	"merge/gitsupports"
	"merge/objects"
	"merge/utils"
	"os"
	"strings"
	"sync"
)

func readInputFile(config objects.Config) (mappingInput map[string]string, keys []string) {
	fmt.Printf(constants.ReadingInputFile, config.InputFile)
	mappingInput, keys = readConfigFile(config.Workspace + constants.PathSeparator + config.InputFile)
	return mappingInput, keys
}

func processInputPath(wg *sync.WaitGroup, path string, config objects.Config, result chan string, partialFileMap map[string]string, skippedFiles map[string]bool, skippedFile chan string) {
	defer wg.Done()
	path = strings.TrimSpace(path)
	_, ok := partialFileMap[path]
	processedPath := skippedFiles[path]
	if processedPath {
		fmt.Printf(constants.FileAlreadyProcessed, path)
		return
	}
	if path != constants.Empty {
		fmt.Printf(constants.ReadingFile, path)

		var contentFromFile string
		if ok {
			fmt.Println("Partial file found, checking for marker")
			var markerIndex int = checkExistingMarker(path, partialFileMap[path])
			if markerIndex > constants.COUNT_MARKER_DEFAULT {
				fmt.Println("Found a marker, reading content from marker to end")
				contentFromFile = readContentFromMarkerToEnd(path, partialFileMap[path], markerIndex)
			} else {
				fmt.Println("Marker not found, reading full content")
				contentFromFile = readContentFile(path)
			}
			appendSignToFile(path, config)
		} else {
			contentFromFile = readContentFile(path)
		}

		if contentFromFile != constants.Empty {
			result <- contentFromFile + constants.BreakLine + config.ConcatChar + constants.BreakLine
		}
		skippedFile <- path
	}
}

func loadPartialFileMap(config objects.Config) (partialFileMap map[string]string) {
	if config.PartialFileMap == constants.Empty {
		return make(map[string]string)
	}
	content := readContentFile(config.Workspace + constants.PathSeparator + config.PartialFileMap)
	partialFileMap = make(map[string]string)
	for line := range strings.SplitSeq(content, constants.BreakLine) {
		line = getAbsolutePath(config.GitRepo + constants.PathSeparator + strings.TrimSpace(line))
		if line != constants.Empty {
			partialFileMap[line] = config.Sign
		}
	}

	return partialFileMap
}

func appendSignToFile(file string, config objects.Config) {
	if config.Sign != constants.Empty {
		copyContent := constants.BreakLine + config.Sign + utils.GetCurrentTime(constants.DateTimeFormat) + constants.BreakLine
		_ = utils.WriteToFile(file, copyContent, os.O_APPEND|os.O_WRONLY)
	}

}

func readPrefixInputFile(config objects.Config) (inputFileSetting map[string]string) {
	if config.PrefixInputFile == constants.Empty {
		return make(map[string]string)

	}
	inputFileSetting, keys := readConfigFile(config.Workspace + constants.PathSeparator + config.PrefixInputFile)
	if inputFileSetting == nil {
		return make(map[string]string)
	}

	for _, key := range keys {
		value := inputFileSetting[key]
		inputFileSetting[key] = getAbsolutePath(config.GitRepo + constants.PathSeparator + value)
	}

	return inputFileSetting
}

func filterPrefixLines(outputFile string, lines []string, prefixMapping map[string]string) []string {
	prefix, hasPrefix := prefixMapping[outputFile]
	filteredLines := make([]string, 0, len(lines))
	if hasPrefix {
		prefixAbs := getAbsolutePath(prefix)
		for _, line := range lines {
			linecopy := getAbsolutePath(line)
			if strings.HasPrefix(linecopy, prefixAbs) {
				filteredLines = append(filteredLines, linecopy)
			}
		}
	}

	return filteredLines
}

func handleLines(lines []string, wg *sync.WaitGroup, config objects.Config, contents chan string, partialFileMap map[string]string, skippedFiles map[string]bool, skippedFile chan string) {
	for _, line := range lines {
		go processInputPath(wg, line, config, contents, partialFileMap, skippedFiles, skippedFile)
	}
}

func processInputFileContent(inputFileSetting map[string]string, keys []string, config objects.Config) {
	prefixMapping := readPrefixInputFile(config)
	partialFileMap := loadPartialFileMap(config)
	skippedFiles := make(map[string]bool, len(inputFileSetting))
	// for outputFile, inputFile := range inputFileSetting {
	for _, outputFile := range keys {

		inputFile := inputFileSetting[outputFile]
		fmt.Printf(constants.ProcessingOutputFile, outputFile)
		outFile := config.Workspace + constants.PathSeparator + config.OutputFolder + constants.PathSeparator + outputFile
		contentFromInput := readContentFile(config.Workspace + constants.PathSeparator + inputFile)
		lines := strings.Split(contentFromInput, constants.BreakLine)
		var contentBuilder strings.Builder
		var wg sync.WaitGroup
		contents := make(chan string, len(lines))
		isConfigPrefix := prefixMapping[outputFile] != constants.Empty
		skippedFile := make(chan string, len(lines))
		if isConfigPrefix {
			filtered := filterPrefixLines(outputFile, lines, prefixMapping)
			// this is very important to add the number of goroutines to wait for
			// if the number is wrong it will cause deadlock
			wg.Add(len(filtered))
			handleLines(filtered, &wg, config, contents, partialFileMap, skippedFiles, skippedFile)
		} else {
			wg.Add(len(lines))
			handleLines(lines, &wg, config, contents, partialFileMap, skippedFiles, skippedFile)
		}
		wg.Wait()          // wait for all goroutines to finish
		close(contents)    // this is important to close the channel after all goroutines are done
		close(skippedFile) // close the skipped file channel
		fmt.Printf("Write content to file: %s \n", outFile)
		for file := range skippedFile {
			skippedFiles[file] = true
		}

		for content := range contents {
			// We don't need to check for empty content here, as we already checked it in processInputPath
			contentBuilder.WriteString(content)
		}

		_ = os.MkdirAll(config.Workspace+constants.PathSeparator+config.OutputFolder, os.ModePerm)
		utils.DeleteFileIfExists(outFile)
		outFileFile, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer outFileFile.Close()
		_, err = outFileFile.WriteString(contentBuilder.String())
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

func readContentFile(filePath string) (content string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf(constants.FileNotFound, filePath)
		return constants.Empty
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var contentBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		contentBuilder.WriteString(line + constants.BreakLine)
	}
	return contentBuilder.String()

}

func checkExistingMarker(path string, marker string) int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(constants.FileNotFound, path)
		return constants.COUNT_MARKER_DEFAULT
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var index int = constants.COUNT_MARKER_DEFAULT
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, marker) {
			index++
		}
	}

	return index
}

func readContentFromMarkerToEnd(path string, marker string, markerIndex int) (result string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(constants.FileNotFound, path)
		return constants.Empty
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var contentBuilder strings.Builder
	foundMarker := false
	countMarker := constants.COUNT_MARKER_DEFAULT
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, marker) {
			countMarker++
			if countMarker >= markerIndex {
				foundMarker = true
			}
			continue
		}
		if foundMarker {
			contentBuilder.WriteString(line + constants.BreakLine)
		}
	}

	return contentBuilder.String()
}

func cleanUp(outputFiles []string, config objects.Config) {
	fmt.Printf(constants.CleaningUpOutputFile)
	for _, outputFile := range outputFiles {
		outFile := getAbsolutePath(config.Workspace + constants.PathSeparator + config.OutputFolder + constants.PathSeparator + outputFile)
		if info, err := os.Stat(outFile); err == nil && info.Size() == 0 {
			utils.DeleteFileIfExists(outFile)
		}
	}
}

func MainLogic() {
	configPath := flag.String("config", constants.Empty, "The config file path.Ex: .mergeConfig")
	checkoutCommit := flag.String("git-show", constants.Empty, "Commit hash to show changed files")
	gitCommand := flag.String("git", constants.Empty, "Git command to execute")
	showConfig := flag.Bool("show-config", false, "Show config values")
	output := flag.String("output", constants.Empty, "Output folder")
	sign := flag.String("sign", constants.Empty, "Sign for the output file, default is SIGNED")
	flag.Parse()
	var config objects.Config
	if *configPath == constants.Empty {
		config = loadConfig(constants.Empty)
	} else {
		config = loadConfig(*configPath)
	}

	if *output != constants.Empty {
		config.OutputFolder = *output
	}

	if *sign != constants.Empty {
		config.Sign = *sign
	}

	if *showConfig {
		fmt.Println("Workspace: ", config.Workspace)
		fmt.Println("Git repo: ", config.GitRepo)
		fmt.Println("Output folder: ", config.OutputFolder)
		fmt.Println("Input file: ", config.InputFile)
		fmt.Println("Sign: ", config.Sign)
		fmt.Println("Concat char: ", config.ConcatChar)
		fmt.Println("Whitelist extensions: ", config.WhitelistExtensions)
	}

	if *checkoutCommit != constants.Empty {
		gitsupports.GetChangedFiles(*checkoutCommit, config)
		return
	}

	if *gitCommand != constants.Empty {
		gitsupports.ExecGitCommand(*gitCommand, config)
		return
	}
	inputMap, keys := readInputFile(config)
	processInputFileContent(inputMap, keys, config)
	cleanUp(keys, config)
}
