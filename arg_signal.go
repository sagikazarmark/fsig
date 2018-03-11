package main

import (
	"strconv"
	"syscall"

	"github.com/alecthomas/kingpin"
)

type signalValue syscall.Signal

func (s *signalValue) Set(value string) error {
	sig, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	*s = (signalValue)(syscall.Signal(sig))

	return nil
}

func (s *signalValue) String() string {
	return s.String()
}

func signalArg(s kingpin.Settings) (target *syscall.Signal) {
	target = new(syscall.Signal)

	s.SetValue((*signalValue)(target))

	return
}
