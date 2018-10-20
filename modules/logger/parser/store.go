// Package logger/parser contains methods to parse and restructure log output from go testing and terratest
package parser

import (
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/files"
	junitformatter "github.com/jstemmer/go-junit-report/formatter"
	junitparser "github.com/jstemmer/go-junit-report/parser"
	"github.com/sirupsen/logrus"
)

type LogWriter struct {
	// Represents an open channel to a test writer for given test indexed by name
	lookup    map[string]chan string
	outputDir string
}

// LogWriter.getOrCreateChannel will get the corresponding channel to a log writer for the provided test name, or create
// a new channel and spawn the corresponding log writer.
func (logWriter LogWriter) getOrCreateChannel(logger *logrus.Logger, testName string) chan<- string {
	writerChan, hasKey := logWriter.lookup[testName]
	if !hasKey {
		writerChan = make(chan string)
		logWriter.lookup[testName] = writerChan
		go collectLogs(logger, logWriter.outputDir, testName, writerChan)
	}
	return writerChan
}

// LogWriter.closeChannels closes all the channels in the lookup dictionary
func (logWriter LogWriter) closeChannels(logger *logrus.Logger) {
	logger.Infof("Closing all the channels in log writer")
	for _, channel := range logWriter.lookup {
		close(channel)
	}
}

// collectLogs will read data off of a channel and write to a log file named as outputDir/testName
// Note that this will drain the channel when we encounter errors, to ensure the upstream functions don't crash.
func collectLogs(logger *logrus.Logger, outputDir string, testName string, writerChan <-chan string) {
	logger.Infof("Spawned log writer for test %s", testName)
	filename := filepath.Join(outputDir, testName+".log")
	logger.Infof("Storing logs for test %s to %s", testName, filename)

	dirName := filepath.Dir(filename)
	err := ensureDirectoryExists(logger, dirName)
	if err != nil {
		logger.Errorf("Error making directory for test %s", testName)
		// Since we don't have a file, simply drain the channel for this log
		drain(writerChan)
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		logger.Errorf("Error making log file for test %s", testName)
		// Since we don't have a file, simply drain the channel for this log
		drain(writerChan)
		return
	}
	defer f.Close()

	for data := range writerChan {
		_, err := f.WriteString(data + "\n")
		if err != nil {
			logger.Errorf("Error (%s) writing log entry: %s", err, data)
		}
	}
	logger.Infof("Channel closed for log writer of test %s", testName)
}

// ensureDirectoryExists will only attempt to create the directory if it does not exist
func ensureDirectoryExists(logger *logrus.Logger, dirName string) error {
	if files.IsDir(dirName) {
		logger.Infof("Directory %s already exists", dirName)
		return nil
	}
	logger.Infof("Creating directory %s", dirName)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		logger.Errorf("Error making directory %s: %s", dirName, err)
	}
	return err
}

// drain simply drains the channel until closed
func drain(channel <-chan string) {
	for _ = range channel {
	}
}

// storeJunitReport takes a parsed Junit report and stores it as report.xml in the output directory
func storeJunitReport(logger *logrus.Logger, outputDir string, report *junitparser.Report) {
	ensureDirectoryExists(logger, outputDir)
	filename := filepath.Join(outputDir, "report.xml")
	f, err := os.Create(filename)
	if err != nil {
		logger.Errorf("Error making file %s for junit report", filename)
		return
	}
	defer f.Close()

	err = junitformatter.JUnitReportXML(report, false, "", f)
	if err != nil {
		logger.Errorf("Error formatting junit xml report: %s", err)
		return
	}
}
