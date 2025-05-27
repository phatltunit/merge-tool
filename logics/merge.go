package logics

import (
	"bufio"
	"flag"
	"fmt"
	"merge/constants"
	"merge/gitsupports"
	"merge/objects"
	"os"
	"strings"
	"sync"
)

func readInputFile(config objects.Config) (mappingInput map[string]string) {
	fmt.Printf(constants.ReadingInputFile, config.InputFile)
	mappingInput = readConfigFile(config.Workspace + constants.PathSeparator + config.InputFile)
	return mappingInput
}

func processInputPath(wg *sync.WaitGroup, path string, config objects.Config, result chan string) {
	defer wg.Done()
	path = strings.TrimSpace(path)
	if path != constants.Empty {
		fmt.Printf(constants.ReadingFile, path)
		contentFromFile := readContentFile(path)
		if contentFromFile != constants.Empty {
			result <- contentFromFile + constants.BreakLine + config.ConcatChar + constants.BreakLine
		}
	}
}

func readPrefixInputFile(config objects.Config) (inputFileSetting map[string]string) {
	if config.PrefixInputFile == constants.Empty {
		return make(map[string]string)

	}
	inputFileSetting = readConfigFile(config.Workspace + constants.PathSeparator + config.PrefixInputFile)
	if inputFileSetting == nil {
		return make(map[string]string)
	}

	for key, value := range inputFileSetting {
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

func handleLines(lines []string, wg *sync.WaitGroup, config objects.Config, contents chan string) {
	for _, line := range lines {
		go processInputPath(wg, line, config, contents)
	}
}

func processInputFileContent(inputFileSetting map[string]string, config objects.Config) {
	prefixMapping := readPrefixInputFile(config)

	for outputFile, inputFile := range inputFileSetting {
		fmt.Printf(constants.ProcessingOutputFile, outputFile)
		outFile := config.Workspace + constants.PathSeparator + config.OutputFolder + constants.PathSeparator + outputFile
		contentFromInput := readContentFile(config.Workspace + constants.PathSeparator + inputFile)
		lines := strings.Split(contentFromInput, constants.BreakLine)
		var contentBuilder strings.Builder
		var wg sync.WaitGroup
		contents := make(chan string, len(lines))
		filtered := filterPrefixLines(outputFile, lines, prefixMapping)
		if len(filtered) > 0 {
			wg.Add(len(filtered))
			handleLines(filtered, &wg, config, contents)
		} else {
			wg.Add(len(lines))
			handleLines(lines, &wg, config, contents)
		}

		wg.Wait()       // wait for all goroutines to finish
		close(contents) // this is important to close the channel after all goroutines are done

		fmt.Printf("Write content to file: %s \n", outFile)
		for content := range contents {
			// We don't need to check for empty content here, as we already checked it in processInputPath
			contentBuilder.WriteString(content)
		}

		_ = os.MkdirAll(config.Workspace+constants.PathSeparator+config.OutputFolder, os.ModePerm)
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

	processInputFileContent(readInputFile(config), config)

}
