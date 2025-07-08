package utils

import (
	"fmt"
	"os/exec"
)

// check n/w connectivity
func CheckNetwork() error {
	cmd := exec.Command("ping", "-c", "1", "8.8.8.8")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Network check failed: %s", string(output))
	}
	return nil
}


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
