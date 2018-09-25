// +build linux darwin

package main

import (
	"os"
	"os/exec"
)

func newChildCommand(cmd string, args []string) *exec.Cmd {
	childCmd := exec.Command(cmd, args...)

	childCmd.Env = os.Environ()
	childCmd.Stdin = os.Stdin
	childCmd.Stdout = os.Stdout
	childCmd.Stderr = os.Stderr

	return childCmd
}

func sendSignal(proc *os.Process, sig os.Signal) error {
	return proc.Signal(sig)
}
