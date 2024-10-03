package tester

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RunChild(testFunc func()) (stdout string, stderr string, err error) {
	// Check if the function is being run as a child process
	if os.Getenv("CHILD_TEST") == "1" {
		testFunc()
		os.Exit(0)
	}

	// Set up the command to run the current process with the test function
	cmd := exec.Command(os.Args[0], "-test.run="+getCallingTestName())
	cmd.Env = append(os.Environ(), "CHILD_TEST=1")

	// Capture the output
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func getCallingTestName() string {
	// Get the name of the calling function
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	funcDetails := runtime.FuncForPC(pc)
	if funcDetails == nil {
		return ""
	}

	name := funcDetails.Name()
	if dotIndex := strings.LastIndex(name, "."); dotIndex >= 0 {
		name = name[dotIndex+1:]
	}
	return name
}
