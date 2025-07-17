package logics

import (
	"bufio"
	"fmt"
	"merge/constants"
	"merge/objects"
	"merge/utils"
	"os"
	"strings"
)

func loadConfig(file string) (config objects.Config) {
	var configMap map[string]string
	if file == constants.Empty {
		configMap, _ = readConfigFile(constants.MergeConfig)
	} else {
		configMap, _ = readConfigFile(file)
	}
	config = objects.Config{
		Workspace:           constants.Empty,
		OutputFolder:        constants.DefailtOutputFolder,
		InputFile:           constants.DefaultInputFile,
		Sign:                constants.DefaultSigned,
		ConcatChar:          constants.DefaultConcatChar,
		WhitelistExtensions: []string{constants.SQL},
		PrefixInputFile:     constants.Empty,
		PartialFileMap:      constants.Empty,
	}
	if configMap != nil {
		config.Workspace = getAbsolutePath(configMap[constants.WorkspaceKey])
		config.OutputFolder = configMap[constants.OutputFolderKey]
		config.InputFile = configMap[constants.InputFileKey]
		config.Sign = configMap[constants.SignKey]
		config.ConcatChar = configMap[constants.ConcatCharKey]
		config.WhitelistExtensions = strings.Split(configMap[constants.WhiteListExtensions], constants.MultipleValuesSeparator)
		config.GitRepo = getAbsolutePath(configMap[constants.GitRepo])
		config.PrefixInputFile = configMap[constants.PrefixInputFile]
		config.PartialFileMap = configMap[constants.PartialFileMap]
	}

	return config
}

func getAbsolutePath(path string) (absPath string) {
	return utils.GetAbsolutePath(path)
}

func readConfigFile(inputFile string) (configMap map[string]string, keys []string) {
	// read the mergeConfig file
	file, err := os.Open(getAbsolutePath(inputFile))
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	configMap = make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		kv := strings.Split(line, constants.PropertiesSeparator)
		if len(kv) != 2 {
			fmt.Printf(constants.InvalidConfigFile, line)
			continue
		}
		key := kv[0]
		value := kv[1]
		configMap[key] = value
		keys = append(keys, key)
	}
	return configMap, keys
}
