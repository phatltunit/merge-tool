# Merge Tool

This tool merges content from multiple files into a single output file based on a configuration.

## Usage

1.  **Configuration:**
    *   Create a `.mergeConfig` file in the root directory of your project, or specify a config file path using the `-config` flag.
    *   The `.mergeConfig` file should contain key-value pairs defining the tool's configuration.

    ```properties
    WORKSPACE=your_workspace_path
    OUTPUT_FOLDER=RESULT
    INPUT_FILE=input.txt
    SIGN=SIGNED
    CONCAT_CHAR=GO
    WHILELIST_EXTENSIONS=.sql
    GIT_REPO=your_git_repo_url
    ```

    *   `WORKSPACE`: The root directory for all input and output files.
    *   `OUTPUT_FOLDER`: The directory where the merged output files will be created (default: `RESULT`).
    *   `INPUT_FILE`: The input file containing a mapping of output files to input files.
    *   `SIGN`: A signature to add to the output files (default: `SIGNED`).
    *   `CONCAT_CHAR`: A character to concatenate the content of the input files (default: `GO`).
    *   `WHILELIST_EXTENSIONS`: A semicolon-separated list of file extensions to include (default: `.sql`).
    *   `GIT_REPO`: The URL of the Git repository (optional).

2.  **Input File:**
    *   Create an `input.txt` file (or specify a different input file in the `.mergeConfig` file) that maps output files to input files. Each line in the input file should be in the format `output_file=input_file`.

    ```properties
    output1.txt=input1.txt
    output2.txt=input2.txt
    ```

3.  **Running the tool:**

    ```bash
    go run main.go
    ```

    *   To specify a config file:

    ```bash
    go run main.go -config your_config_file.config
    ```

    *   To get changed files from a git commit:

    ```bash
    go run main.go -git-show commit_hash
    ```

## Example

1.  Create a `.mergeConfig` file:

    ```properties
    WORKSPACE=d:/Backend/merge
    OUTPUT_FOLDER=RESULT
    INPUT_FILE=input.txt
    SIGN=SIGNED
    CONCAT_CHAR=GO
    WHILELIST_EXTENSIONS=.sql
    ```

2.  Create an `input.txt` file:

    ```properties
    output.sql=file1.sql
    ```

3.  Create `file1.sql` in the `d:/Backend/merge` directory:

    ```sql
    SELECT * FROM table1;
    ```

4.  Run the tool:

    ```bash
    go run main.go
    ```

5.  The merged content will be in `d:/Backend/merge/RESULT/output.sql`.

## Git Support

The tool can also be used to get changed files from a git commit. To do this, use the `-git-show` flag:

```bash
go run main.go -git-show commit_hash
```

This will print the list of changed files in the specified commit. To specify multiple commit hashes, separate them with commas:

```bash
go run main.go -git-show commit_hash1,commit_hash2,commit_hash3
```

This will print the list of changed files across all specified commits. The files will be distinct and ordered by alphabet, which will help in some cases.

## Building an Executable

To build this project into a single executable file for your machine, use the following command:

```bash
go build main.go
```

This will create an executable file named `main.exe` on Window or `main` on Linux in the root directory of the project. You can then run this executable directly.

## Note
Make sure you have Go installed and configured correctly.
