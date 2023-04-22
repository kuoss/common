package tester_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestCaptureChildTest_noExit(t *testing.T) {
	stdout, stderr, err := tester.CaptureChildTest(func() {
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
	stdout, stderr, err := tester.CaptureChildTest(func() {
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
	stdout, stderr, err := tester.CaptureChildTest(func() {
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
