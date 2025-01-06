package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Blacklist []string `yaml:"blacklist"`
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config.yaml")
}

func loadConfig() Config {
	configPath := getConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("YAML file not found. Creating default config at %s\n", configPath)
		defaultConfig := Config{
			Blacklist: []string{
				"otf", "woff2", "js", "ttf", "woff", "eot",
				"svg", "png", "jpg", "jpeg", "gif", "bmp",
				"css", "webp", "tiff", "heic", "heif",
			},
		}

		file, err := os.Create(configPath)
		if err != nil {
			log.Fatalf("Failed to create YAML file: %v", err)
		}
		defer file.Close()

		encoder := yaml.NewEncoder(file)
		if err := encoder.Encode(defaultConfig); err != nil {
			log.Fatalf("Failed to write default config to YAML file: %v", err)
		}
		log.Printf("Default config written to %s\n", configPath)
		return defaultConfig
	}

	// Load config from YAML file
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error loading YAML file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	return config
}

func filterURLs(urls []string, blacklist []string) []string {
	var cleanURLs []string
	for _, url := range urls {
		// Remove query parameters to get the file part of the URL
		re := regexp.MustCompile(`(\?[^\s]*)|(#.*)$`)
		urlWithoutParams := re.ReplaceAllString(url, "")

		// Extract file extension (handling case insensitivity)
		ext := strings.ToLower(filepath.Ext(urlWithoutParams))

		// Check if the extension is in the blacklist
		excluded := false
		for _, pattern := range blacklist {
			if ext == "."+strings.ToLower(pattern) {
				excluded = true
				break
			}
		}
		if !excluded {
			cleanURLs = append(cleanURLs, url)
		}
	}
	return cleanURLs
}

func readURLsFromInput() []string {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
	return urls
}

func readURLsFromFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from file: %v", err)
	}
	return urls
}

func saveURLsToFile(urls []string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	for _, url := range urls {
		_, err := file.WriteString(url + "\n")
		if err != nil {
			log.Fatalf("Error writing to output file: %v", err)
		}
	}
}

func printUsage() {
	fmt.Println("Usage: ./main [OPTIONS]")
	fmt.Println("A program to filter URLs based on a blacklist from a YAML configuration file.")
	fmt.Println("Author : sudomode | LinkedIn : https://www.linkedin.com/in/0xsudomode ")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -i <file>    Input file containing URLs (one per line)")
	fmt.Println("  -o <file>    Output file to save filtered URLs")
	fmt.Println("")
	fmt.Println("If no input or output file is provided, the program reads from stdin and writes to stdout.")
}

func isInputFromPipe() bool {
	// Check if there's input from stdin (a pipe is used)
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Error checking stdin: %v", err)
	}
	// If the file mode is named pipe or has data, it's a pipe.
	return fi.Mode()&os.ModeNamedPipe != 0
}

func main() {
	// Load config if it exists or create it if it doesn't
	config := loadConfig()

	// Check if the program is run without parameters and is not piped input
	if len(os.Args) == 1 && !isInputFromPipe() {
		// Display usage and exit if no parameters are provided
		printUsage()
		os.Exit(1)
	}

	var urls []string
	inputFile := ""
	outputFile := ""

	// Check for input and output options
	for i, arg := range os.Args {
		if arg == "-i" && i+1 < len(os.Args) {
			inputFile = os.Args[i+1]
		} else if arg == "-o" && i+1 < len(os.Args) {
			outputFile = os.Args[i+1]
		}
	}

	if inputFile != "" {
		urls = readURLsFromFile(inputFile)
	} else {
		urls = readURLsFromInput()
	}

	cleanURLs := filterURLs(urls, config.Blacklist)

	if outputFile != "" {
		saveURLsToFile(cleanURLs, outputFile)
	} else {
		for _, url := range cleanURLs {
			fmt.Println(url)
		}
	}
}
