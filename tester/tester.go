package tester

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CaptureChildTest(f func()) (stdout string, stderr string, err error) {
	if os.Getenv("CHILD") == "1" {
		f()
		return "", "", nil
	}
	testRun, err := getTestRun()
	if err != nil {
		return "", "", fmt.Errorf("error on getTestRun: %w", err)
	}
	cmd := exec.Command(os.Args[0], "-test.run="+testRun)
	cmd.Env = append(os.Environ(), "CHILD=1")
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("error on Run: %w", err)
	}
	return stdoutBuf.String(), stderrBuf.String(), nil
}

func getTestRun() (string, error) {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "", fmt.Errorf("not ok on Caller")
	}
	name := runtime.FuncForPC(pc).Name()
	if dot := strings.LastIndex(name, "."); dot >= 0 {
		name = name[dot+1:]
	}
	return name, nil
}
