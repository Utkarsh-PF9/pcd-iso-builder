package utils

import (
	"fmt"
	"os/exec"
)

// Helper function to run a command
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = stdOutLogger{}
	cmd.Stderr = stdErrLogger{}
	return cmd.Run()
}

type stdOutLogger struct{}
type stdErrLogger struct{}

func (stdOutLogger) Write(p []byte) (int, error) {
	return fmt.Print("[stdout] ", string(p))
}

func (stdErrLogger) Write(p []byte) (int, error) {
	return fmt.Print("[stderr] ", string(p))
}
