package tester

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CaptureChildTest(f func()) (stdout string, stderr string, err error) {
	if os.Getenv("CHILD") == "1" {
		f()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run="+getTestRun())
	cmd.Env = append(os.Environ(), "CHILD=1")
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func getTestRun() string {
	pc, _, _, _ := runtime.Caller(2) // always ok
	name := runtime.FuncForPC(pc).Name()
	if dot := strings.LastIndex(name, "."); dot >= 0 {
		name = name[dot+1:]
	}
	return name
}
