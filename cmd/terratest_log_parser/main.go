package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/logger/parser"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Error finding current directory: %s", err)
	}
	defaultOutputDir := filepath.Join(currentDir, "out")

	filename := flag.String("testlog", "", "Path to file containing test log. If unset will use stdin.")
	outputDirectory := flag.String(
		"outputdir", defaultOutputDir, "Path to directory to output test output to. If unset will use the current directory.")
	flag.Parse()

	var file *os.File
	if *filename != "" {
		logger.Infof("reading from file")
		file, err = os.Open(*filename)
		if err != nil {
			logger.Fatalf("Error opening file: %s", err)
		}
	} else {
		logger.Infof("reading from stdin")
		file = os.Stdin
	}
	defer file.Close()

	*outputDirectory, err = filepath.Abs(*outputDirectory)
	if err != nil {
		logger.Fatalf("Error extracting absolute path of output directory: %s", err)
	}

	parser.SpawnParsers(logger, file, *outputDirectory)
}
