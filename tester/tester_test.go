package tester

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureChildTest(t *testing.T) {
	stdout, stderr, err := CaptureChildTest(func() {
		fmt.Print("hello")
		print("world")
		os.Exit(2)
	})
	assert.Equal(t, "hello", stdout)
	assert.Equal(t, "world", stderr)
	assert.Error(t, err, "exit status 2")
}

func TestGetTestRun(t *testing.T) {
	assert.Equal(t, "tRunner", getTestRun())
	// always ok at any depth
	funcA := func() string {
		return getTestRun()
	}
	assert.Equal(t, "TestGetTestRun", funcA())
	funcB := func() string {
		return funcA()
	}
	assert.Equal(t, "func2", funcB())
	funcC := func() string {
		return funcA()
	}
	assert.Equal(t, "func3", funcC())
}
