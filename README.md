# Merge Tool

## 1. Configuration (`.mergeConfig` file)

The `.mergeConfig` file contains the main configuration parameters for the merge process:

*   **WORKSPACE**: The location containing configuration files. The `.mergeConfig` file does not necessarily need to be in this directory.
*   **WHILELIST_EXTENSIONS**: File extensions that will be merged (e.g., `.sql`, `.txt`).
*   **CONCAT_CHAR**: The character or string used to concatenate content when merging multiple SQL files (e.g., `GO` or `\`).
*   **INPUT_FILE**: A file that defines the mapping between input files and their corresponding output files.
*   **PREFIX_INPUT_FILE**: Defines paths to files with prefixes that will be merged into the corresponding output. Typically `input_prefix.txt`.
*   **PARTIAL_FILE_MAP**: Specifies files from which only a partial content should be taken, starting from the position indicated by the `SIGN` configuration value until the end of the file.
*   **SIGN**: A unique value automatically appended to the end of the file when processing partial files. This value is used to automatically find the starting point for subsequent processing.
*   **GIT_REPO**: The root directory containing the `.git` directory.
*   **OUTPUT_FOLDER**: Defines the output files relative to the `WORKSPACE` directory.

**Note**: The `GIT_REPO`, `WORKSPACE`, `CONCAT_CHAR`, `PREFIX_INPUT_FILE`, and `PARTIAL_FILE_MAP` configurations need to be redefined based on usage.

## 2. Usage

The program reads the `INPUT_FILE` to determine which configuration file to read for each output file. Currently, all configurations will be read from `all_input.txt` for convenience, allowing for a single copy-paste operation (e.g., from a `git show` command). Therefore, another configuration is needed to map based on prefixes using `PREFIX_INPUT_FILE`.

With this configuration, the program will read from the `input_prefix.txt` file. In this file, each prefix path will correspond to an output.

This handles cases where multiple files are concatenated into a single output.

For cases where only a partial content from a file is needed, starting from the `SIGN` value, another configuration, `PARTIAL_FILE_MAP`, will be added. This file specifies which files should only have partial content taken.

There are cases where a file from which only partial content is needed is within the same prefix as other outputs. Therefore, we will prioritize files with the closest prefix to the file (refer to the sample configuration).

## 3. Other Features

Additionally, the program supports Git:

*   **`--git-show <commit_hash>`**: Retrieves a list of files in that commit. Multiple commits can be retrieved simultaneously by separating each `<hash_commit>` with a comma (`,`) without spaces.

To use `main.exe` from the command line, add its directory to the `PATH` environment variable and optionally rename the file (e.g., to `merge.exe`).

*   **`merge --git-show <hash1>,<hash2>`**: Retrieves files from two commits, `<hash1>` and `<hash2>`.
*   **`merge --help`**: Displays help instructions.
*   **`merge --config <path_to_mergeconfig>`**: Performs a merge according to the configuration from the specified config file, callable from anywhere.
*   **`merge --output <path>`**: Overrides the `OUTPUT_FOLDER` configuration.
*   **`merge --sign <sign>`**: Overrides the `SIGN` configuration.
*   **`merge --git <git_command>`**: Runs a Git command.
*   **`merge --show-config`**: Displays the loaded configuration.
