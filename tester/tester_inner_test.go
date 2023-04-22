package tester

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureChildTest_noExit(t *testing.T) {
	stdout, stderr, err := CaptureChildTest(func() {
		fmt.Fprintf(os.Stdout, "Hello")
		fmt.Fprintf(os.Stderr, "Lorem")
		fmt.Fprintf(os.Stdout, "World")
		fmt.Fprintf(os.Stderr, "Ipsum")
	})
	assert.Equal(t, "HelloWorld", stdout)
	assert.Equal(t, "LoremIpsum", stderr)
	assert.Nil(t, err)
}

func TestCaptureChildTest_exit0(t *testing.T) {
	stdout, stderr, err := CaptureChildTest(func() {
		fmt.Fprintf(os.Stdout, "Hello")
		fmt.Fprintf(os.Stderr, "Lorem")
		fmt.Fprintf(os.Stdout, "World")
		fmt.Fprintf(os.Stderr, "Ipsum")
		os.Exit(0)
	})
	assert.Equal(t, "HelloWorld", stdout)
	assert.Equal(t, "LoremIpsum", stderr)
	assert.Nil(t, err)
}

func TestCaptureChildTest_exit2(t *testing.T) {
	stdout, stderr, err := CaptureChildTest(func() {
		fmt.Fprintf(os.Stdout, "Hello")
		fmt.Fprintf(os.Stderr, "Lorem")
		fmt.Fprintf(os.Stdout, "World")
		fmt.Fprintf(os.Stderr, "Ipsum")
		os.Exit(2)
	})
	assert.Equal(t, "HelloWorld", stdout)
	assert.Equal(t, "LoremIpsum", stderr)
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
