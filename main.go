package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kvnxiao/sort-awesome-lists/logging"
	"github.com/kvnxiao/sort-awesome-lists/parser"
)

func main() {
	// flags setup
	tokenPtr := flag.String("t", "", "GitHub personal access token")
	verbosePtr := flag.Bool("v", false, "prints debug messages to stdout if true (default = false)")
	outputPtr := flag.String("o", "", "name of file to write output to if set, otherwise prints to stdout")
	subBlockSizePtr := flag.Int("bs", 10, "number of concurrent requests to send to GitHub API at a time, per each block found (default = 10).")
	flag.Parse()

	// read token
	token := *tokenPtr
	if token == "" {
		log.Fatalf("Please pass in a GitHub personal access token before using")
	}

	// set verbosity in logger
	verbose := *verbosePtr
	logging.SetVerbose(verbose)

	// parse args for link
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("A URL to the markdown file must be provided!")
	}
	link := args[0]
	logging.Printlnf("URL to parse markdown: %s", link)

	// check file path
	outputFileName := *outputPtr
	outputFilePath := checkFilePath(outputFileName)

	// parse and sort markdown by number of github stars
	md := parseAndSort(link, token, *subBlockSizePtr)
	sortedContents := md.ToString()

	if outputFilePath != "" {
		err := ioutil.WriteFile(outputFileName, []byte(sortedContents), 0666)
		if err != nil {
			log.Fatalf("failed to write to file %s: %v", outputFileName, err)
		}
	} else {
		fmt.Println(sortedContents)
	}
}

func checkFilePath(path string) string {
	if path == "" {
		return ""
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("specified output path is invalid: %s", absPath)
	}
	if fileExists(absPath) {
		log.Fatalf("file already exists in path %s", absPath)
	}
	return absPath
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func parseAndSort(link, token string, subBlockSize int) *parser.Markdown {
	md := parser.ParseMarkdown(link)
	md.FetchStars(token, subBlockSize)
	md.Sort()
	return md
}
