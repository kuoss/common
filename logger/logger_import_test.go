package logger_test

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/kuoss/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestSetLevel(t *testing.T) {
	for _, level := range logger.AllLevels {
		logger.SetLevel(level)
		assert.Equal(t, level, logger.GetLevel())
	}

	logger.SetLevel(logger.DebugLevel)
	output1 := captureOutput(func() {
		logger.Debugf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=debug msg="hello=world lorem=ipsum number=42" file="logger_import_test.go:[0-9]+"`, output1)

	logger.SetLevel(logger.InfoLevel)
	output2 := captureOutput(func() {
		logger.Debugf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	assert.Equal(t, "", output2)
}

func TestSetFullpath(t *testing.T) {

	logger.SetFullpath(true)
	output := captureOutput(func() {
		logger.Warnf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world lorem=ipsum number=42" file="[^"]+/common/logger/logger_import_test.go:[0-9]+"`, output)

	logger.SetFullpath(false)
	output = captureOutput(func() {
		logger.Warnf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world lorem=ipsum number=42" file="logger_import_test.go:[0-9]+"`, output)
}

func captureOutput(f func()) string {
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)
	f()
	logger.SetOutput(os.Stderr)
	return buf.String()
}

func TestDebugf(t *testing.T) {
	output := captureOutput(func() {
		logger.Debugf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	assert.Equal(t, "", output)
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		logger.Infof("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=info msg="hello=world lorem=ipsum number=42" file="logger_import_test.go:[0-9]+"`, output)
}

func TestWarnf(t *testing.T) {
	output := captureOutput(func() {
		logger.Warnf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world lorem=ipsum number=42" file="logger_import_test.go:[0-9]+"`, output)
}

func TestErrorf(t *testing.T) {
	output := captureOutput(func() {
		logger.Errorf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=error msg="hello=world lorem=ipsum number=42" file="logger_import_test.go:[0-9]+"`, output)
}

func TestFatalf_import(t *testing.T) {
	if os.Getenv("FATALF") == "1" {
		logger.Fatalf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalf_import")
	cmd.Env = append(os.Environ(), "FATALF=1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := stderr.String()
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=fatal msg="hello=world lorem=ipsum number=42" file="logger_import_test.go:[0-9]+"`, output)
	assert.Error(t, err, "exit status 1")
}
