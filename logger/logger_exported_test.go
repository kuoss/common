package logger_test

import (
	"bytes"
	"github.com/kuoss/common/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestSettingFullpath(t *testing.T) {
	logger.SetFullpath(true)
	output := captureOutput(func() {
		logger.Warnf("hello=%s lorem=%s number=%d", "hello", "ipsum", 42)
	})

	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=hello lorem=ipsum number=42" file="[^"]+/common/logger/logger_exported_test.go:[0-9]+"`, output)

	logger.SetFullpath(false)
	output = captureOutput(func() {
		logger.Warnf("hello=%s lorem=%s number=%d", "hello", "ipsum", 42)
	})

	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=hello lorem=ipsum number=42" file="logger_exported_test.go:[0-9]+"`, output)
}

func captureOutput(f func()) string {
	buf := &bytes.Buffer{}
	logger.SetOutput(io.MultiWriter(buf, os.Stderr))
	f()
	return buf.String()
}
