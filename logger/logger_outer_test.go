package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestSetLevel(t *testing.T) {
	for _, level := range logger.AllLevels {
		logger.SetLevel(level)
		assert.Equal(t, level, logger.GetLevel())
	}

	logger.SetLevel(logger.DebugLevel)
	output1 := captureOutput(func() {
		logger.Debugf("hello=%s number=%d", "world", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=debug msg="hello=world number=42" file="logger_outer_test.go:[0-9]+"`, output1)

	logger.SetLevel(logger.InfoLevel)
	output2 := captureOutput(func() {
		logger.Debugf("hello=%s number=%d", "world", 42)
	})
	assert.Equal(t, "", output2)
}

func TestSetFullpath(t *testing.T) {

	logger.SetFullpath(true)
	output := captureOutput(func() {
		logger.Warnf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world number=42" file="[^"]+/common/logger/logger_outer_test.go:[0-9]+"`, output)

	logger.SetFullpath(false)
	output = captureOutput(func() {
		logger.Warnf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world number=42" file="logger_outer_test.go:[0-9]+"`, output)
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
		logger.Debugf("hello=%s number=%d", "world", 42)
	})
	assert.Equal(t, "", output)
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		logger.Infof("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=info msg="hello=world number=42" file="logger_outer_test.go:[0-9]+"`, output)
}

func TestWarnf(t *testing.T) {
	output := captureOutput(func() {
		logger.Warnf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world number=42" file="logger_outer_test.go:[0-9]+"`, output)
}

func TestErrorf(t *testing.T) {
	output := captureOutput(func() {
		logger.Errorf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=error msg="hello=world number=42" file="logger_outer_test.go:[0-9]+"`, output)
}

func TestFatalf_outer(t *testing.T) {
	_, output, err := tester.CaptureChildTest(func() {
		logger.Fatalf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=fatal msg="hello=world number=42" file="logger_outer_test.go:[0-9]+"`, output)
	assert.Error(t, err, "exit status 1")
}
