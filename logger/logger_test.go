package logger

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	dummyFrame = &runtime.Frame{Function: "github.com/kuoss/example/pkg.Func1", File: "/example/pkg/file1.go"}
)

func TestInit(t *testing.T) {
	assert.NotNil(t, logger)
	assert.NotZero(t, logger)
	assert.Equal(t, logger.Out, os.Stderr)

	assert.NotZero(t, logger.Formatter)
	assert.True(t, logger.Formatter.(*logrus.TextFormatter).FullTimestamp)
	assert.NotZero(t, logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier)

	funcname, filename := logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier(dummyFrame)
	assert.Equal(t, "", funcname)
	assert.Equal(t, "???:1", filename)
}

func TestGetCallerPrettyfier(t *testing.T) {
	var funcname string
	var filename string

	funcname, filename = getCallerPrettyfier(false)(dummyFrame)
	assert.Equal(t, "", funcname)
	assert.Equal(t, "???:1", filename)

	funcname, filename = getCallerPrettyfier(true)(dummyFrame)
	assert.Equal(t, "", funcname)
	assert.Equal(t, "???:1", filename)
}

func TestSetLevel(t *testing.T) {
	for _, level := range AllLevels {
		SetLevel(level)
		assert.Equal(t, level, GetLevel())
	}

	SetLevel(DebugLevel)
	output1 := captureOutput(func() {
		Debugf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=debug msg="hello=world lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output1)

	SetLevel(InfoLevel)
	output2 := captureOutput(func() {
		Debugf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	assert.Equal(t, "", output2)
}

func TestSetFullpath(t *testing.T) {
	SetFullpath(true)
	output := captureOutput(func() {
		Warnf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world lorem=ipsum number=42" file="[^"]+/common/logger/logger_test.go:[0-9]+"`, output)
	SetFullpath(false)
	output = captureOutput(func() {
		Warnf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}

func captureOutput(f func()) string {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	f()
	SetOutput(os.Stderr)
	return buf.String()
}

func TestDebugf(t *testing.T) {
	output := captureOutput(func() {
		Debugf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	assert.Equal(t, "", output)
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		Infof("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=info msg="hello=world lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}

func TestWarnf(t *testing.T) {
	output := captureOutput(func() {
		Warnf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}

func TestErrorf(t *testing.T) {
	output := captureOutput(func() {
		Errorf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=error msg="hello=world lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}

func TestFatalf(t *testing.T) {
	if os.Getenv("FATALF") == "1" {
		Fatalf("hello=%s lorem=%s number=%d", "world", "ipsum", 42)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalf")
	cmd.Env = append(os.Environ(), "FATALF=1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := stderr.String()
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=fatal msg="hello=world lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
	assert.Error(t, err, "exit status 1")
}
