// +build !linux,!amd64 !darwin,!amd64

package main

var signals map[string]syscall.Signal
