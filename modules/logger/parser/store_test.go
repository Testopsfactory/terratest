package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func createLogWriter(t *testing.T) LogWriter {
	dir := getTempDir(t)
	logWriter := LogWriter{
		lookup:    make(map[string]chan string),
		outputDir: dir,
	}
	return logWriter
}

func TestEnsureDirectoryExistsCreatesDirectory(t *testing.T) {
	t.Parallel()

	dir := getTempDir(t)
	defer os.RemoveAll(dir)

	logger := NewTestLogger(t)
	tmpd := filepath.Join(dir, "tmpdir")
	assert.False(t, files.IsDir(tmpd))
	ensureDirectoryExists(logger, tmpd)
	assert.True(t, files.IsDir(tmpd))
}

func TestEnsureDirectoryExistsHandlesExistingDirectory(t *testing.T) {
	t.Parallel()

	dir := getTempDir(t)
	defer os.RemoveAll(dir)

	logger := NewTestLogger(t)
	assert.True(t, files.IsDir(dir))
	ensureDirectoryExists(logger, dir)
	assert.True(t, files.IsDir(dir))
}

func TestGetOrCreateChannelCreatesNewChannel(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	channel := logWriter.getOrCreateChannel(logger, "TestGetOrCreateChannelCreatesNewChannel")
	defer close(channel)
	assert.NotNil(t, channel)
}

func TestGetOrCreateChannelReturnsExistingChannel(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	testName := t.Name()
	logger := NewTestLogger(t)
	var writeOnlyChannel chan<- string
	channel := make(chan string)
	writeOnlyChannel = channel
	defer close(channel)
	logWriter.lookup[testName] = channel
	lookupChannel := logWriter.getOrCreateChannel(logger, testName)
	assert.Equal(t, lookupChannel, writeOnlyChannel)
}

func TestLogCollectorCreatesAndWritesToFile(t *testing.T) {
	t.Parallel()

	testName := t.Name()
	dir := getTempDir(t)
	defer os.RemoveAll(dir)

	logger := NewTestLogger(t)
	channel := make(chan string)

	var waitForCollector sync.WaitGroup
	waitForCollector.Add(1)
	go func() {
		defer waitForCollector.Done()
		collectLogs(logger, dir, testName, channel)
	}()

	randomString := random.UniqueId()
	channel <- randomString
	close(channel)

	// give time for logcollector to finish
	time.Sleep(1 * time.Second)

	logFileName := filepath.Join(dir, testName+".log")
	content, err := ioutil.ReadFile(logFileName)
	assert.Nil(t, err)
	assert.Equal(t, string(content), randomString+"\n")
}

func TestGetOrCreateChannelSpawnsLogCollectorOnCreate(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	testName := t.Name()
	channel := logWriter.getOrCreateChannel(logger, testName)
	assert.NotNil(t, channel)

	randomString := random.UniqueId()
	channel <- randomString
	close(channel)

	// give time for logcollector to finish
	time.Sleep(1 * time.Second)

	logFileName := filepath.Join(logWriter.outputDir, testName+".log")
	content, err := ioutil.ReadFile(logFileName)
	assert.Nil(t, err)
	assert.Equal(t, string(content), randomString+"\n")
}

func TestCloseChannelsClosesAll(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	testName := t.Name()
	alternativeTestName := t.Name() + "Alternative"
	logger := NewTestLogger(t)
	channel1 := make(chan string)
	channel2 := make(chan string)
	logWriter.lookup[testName] = channel1
	logWriter.lookup[alternativeTestName] = channel2

	var waitForClosedChannels sync.WaitGroup
	waitForClosedChannels.Add(2)
	go func() {
		<-channel1
		waitForClosedChannels.Done()
	}()
	go func() {
		<-channel2
		waitForClosedChannels.Done()
	}()
	logWriter.closeChannels(logger)
	waitForClosedChannels.Wait()
}
