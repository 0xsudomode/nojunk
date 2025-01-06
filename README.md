# nojunk

This Go program filters URLs based on a blacklist of file extensions, removing any URLs that end with extensions specified in a YAML configuration file. It also handles URL parameters correctly, ensuring that files with blacklisted extensions (even if they have parameters) are excluded from the output.

## Features

- Reads URLs from input (stdin or a file).
- Filters out URLs that match blacklisted extensions (from a YAML config).
- Handles URLs with query parameters (e.g., `file.js?v=1628142064`).
- Outputs the filtered URLs to stdout or saves them to a file.
- Configurable blacklist stored in a YAML file.

## Requirements

- Go 1.18 or later
- YAML library for Go (`gopkg.in/yaml.v2`)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/nojunk.git
   cd nojunk
   ```
   
2. Install dependencies:

  ```bash
   go get gopkg.in/yaml.v2
  ```
3. Build the tool :

  ```bash
  go build nojunk.go
  mv nojunk ~/go/bin/ # optional
  ```

## Usage

### Configuration

The program uses a configuration file (`.config.yaml`) stored in the user's home directory to specify the file extensions to be blacklisted. If this file doesn't exist, the program will create a default configuration with commonly blacklisted extensions (e.g., `.jpg`, `.css`, `.js`, etc.).

### Command-line Options

-   `-i <file>`: Specify an input file containing URLs (one per line).
-   `-o <file>`: Specify an output file to save filtered URLs.
-   If no input or output files are provided, the program will read URLs from stdin and output to stdout.

### Example

1.  **Filtering URLs from stdin**:
    
   You can pipe URLs to the program, and it will output the filtered URLs:
    
  ```bash
    echo "https://example.com" | gau | nojunk
  ```
2.  **Filtering URLs from a file**:

   ```bash
    nojunk -i input.txt -o output.txt
  ```

If the configuration file (.config.yaml) doesn't exist, the program will create a default one. The blacklist will include extensions like .jpg, .css, .png, etc.


