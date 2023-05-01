package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	dummyFrame = &runtime.Frame{Function: "github.com/kuoss/example/pkg.Func1", File: "/example/pkg/file1.go"}
)

func init() {
	SetCallerSkip(9) // for go test
}

func TestInit(t *testing.T) {
	assert.NotEmpty(t, logger)
	assert.Equal(t, logger.Out, os.Stderr)

	assert.NotEmpty(t, logger.Formatter)
	assert.True(t, logger.Formatter.(*logrus.TextFormatter).FullTimestamp)
	assert.NotEmpty(t, logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier)

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
		Debugf("hello=%s number=%d", "world", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=debug msg="hello=world number=42" file="logger_inner_test.go:[0-9]+"`, output1)

	SetLevel(InfoLevel)
	output2 := captureOutput(func() {
		Debugf("hello=%s number=%d", "world", 42)
	})
	assert.Equal(t, "", output2)
}

func TestSetFullpath(t *testing.T) {
	SetFullpath(true)
	output := captureOutput(func() {
		Warnf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world number=42" file="[^"]+/common/logger/logger_inner_test.go:[0-9]+"`, output)

	SetFullpath(false)
	output = captureOutput(func() {
		Warnf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world number=42" file="logger_inner_test.go:[0-9]+"`, output)
}

func TestSetCallerSkip(t *testing.T) {
	testCases := []struct {
		skip         int
		wantContains string
	}{
		{0, `level=warning msg="hello=world number=42" file="logger.go:74"`},
		{1, `level=warning msg="hello=world number=42" file="text_formatter.go:159"`},
		{2, `level=warning msg="hello=world number=42" file="entry.go:289"`},
		{3, `level=warning msg="hello=world number=42" file="entry.go:252"`},
		{4, `level=warning msg="hello=world number=42" file="entry.go:304"`},
		{5, `level=warning msg="hello=world number=42" file="entry.go:349"`},
		{6, `level=warning msg="hello=world number=42" file="logger.go:154"`},
		{7, `level=warning msg="hello=world number=42" file="logger.go:178"`},
		{8, `level=warning msg="hello=world number=42" file="logger.go:98"`},
		{9, `level=warning msg="hello=world number=42" file="logger_inner_test.go:110"`}, // good for go test
		{10, `level=warning msg="hello=world number=42" file="logger_inner_test.go:121"`},
		{11, `level=warning msg="hello=world number=42" file="logger_inner_test.go:109"`},
		{12, `level=warning msg="hello=world number=42" file="testing.go:1576"`},
		{13, `level=warning msg="hello=world number=42" file="asm_amd64.s:1598"`},
		{14, `level=warning msg="hello=world number=42" file="???:1"`},
		{15, `level=warning msg="hello=world number=42" file="???:1"`},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("skip=%d", tc.skip), func(t *testing.T) {
			SetCallerSkip(tc.skip)
			got := captureOutput(func() {
				Warnf("hello=%s number=%d", "world", 42)
			})
			assert.Contains(t, got, tc.wantContains)
		})
	}
	SetCallerSkip(9)
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
		Debugf("hello=%s number=%d", "world", 42)
	})
	assert.Equal(t, "", output)
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		Infof("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=info msg="hello=world number=42" file="logger_inner_test.go:[0-9]+"`, output)
}

func TestWarnf(t *testing.T) {
	output := captureOutput(func() {
		Warnf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=world number=42" file="logger_inner_test.go:[0-9]+"`, output)
}

func TestErrorf(t *testing.T) {
	output := captureOutput(func() {
		Errorf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=error msg="hello=world number=42" file="logger_inner_test.go:[0-9]+"`, output)
}

func TestFatalf(t *testing.T) {
	_, output, err := tester.CaptureChildTest(func() {
		Fatalf("hello=%s number=%d", "world", 42)
	})
	t.Log(output)
	assert.Regexp(t, `time="[^"]+" level=fatal msg="hello=world number=42" file="logger_inner_test.go:[0-9]+"`, output)
	assert.Error(t, err, "exit status 1")
}
