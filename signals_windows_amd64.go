package main

import (
	"os"
	"os/exec"
	"syscall"
)

var signals = map[string]syscall.Signal{
	"INT":  syscall.SIGINT,
	"KILL": syscall.SIGKILL,
	"QUIT": syscall.SIGQUIT,

	"SIGINT":  syscall.SIGINT,
	"SIGKILL": syscall.SIGKILL,
	"SIGQUIT": syscall.SIGQUIT,
}

func newChildCommand(cmd string, args []string) *exec.Cmd {
	childCmd := exec.Command(cmd, args...)

	childCmd.Env = os.Environ()
	childCmd.Stdin = os.Stdin
	childCmd.Stdout = os.Stdout
	childCmd.Stderr = os.Stderr
	childCmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT | 0x00000200,
	}

	return childCmd
}

func sendSignal(proc *os.Process, sig os.Signal) error {
	if sig == syscall.SIGKILL {
		return proc.Signal(sig)
	}

	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	f, err := dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}
	r1, _, err := f.Call(0, 1)
	if r1 == 0 {
		return err
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}

	event := syscall.CTRL_C_EVENT
	if sig == syscall.SIGQUIT {
		event = syscall.CTRL_BREAK_EVENT
	}
	r1, _, err = f.Call(uintptr(event), uintptr(proc.Pid))
	if r1 == 0 {
		return err
	}
	return nil
}
