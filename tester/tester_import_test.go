package tester_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestCaptureChildTest(t *testing.T) {
	stdout, stderr, err := tester.CaptureChildTest(func() {
		fmt.Print("hello")
		print("world")
		os.Exit(2)
	})
	assert.Equal(t, "hello", stdout)
	assert.Equal(t, "world", stderr)
	assert.Error(t, err, "exit status 2")
}
