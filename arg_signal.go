package main

import (
	"fmt"
	"syscall"

	"github.com/alecthomas/kingpin"
)

var signalTable = map[string]syscall.Signal{
	// With SIG prefix
	"SIGABRT":   syscall.SIGABRT,
	"SIGALRM":   syscall.SIGALRM,
	"SIGBUS":    syscall.SIGBUS,
	"SIGCHLD":   syscall.SIGCHLD,
	"SIGCONT":   syscall.SIGCONT,
	"SIGEMT":    syscall.SIGEMT,
	"SIGFPE":    syscall.SIGFPE,
	"SIGHUP":    syscall.SIGHUP,
	"SIGILL":    syscall.SIGILL,
	"SIGINFO":   syscall.SIGINFO,
	"SIGINT":    syscall.SIGINT,
	"SIGIO":     syscall.SIGIO,
	"SIGIOT":    syscall.SIGIOT,
	"SIGKILL":   syscall.SIGKILL,
	"SIGPIPE":   syscall.SIGPIPE,
	"SIGPROF":   syscall.SIGPROF,
	"SIGQUIT":   syscall.SIGQUIT,
	"SIGSEGV":   syscall.SIGSEGV,
	"SIGSTOP":   syscall.SIGSTOP,
	"SIGSYS":    syscall.SIGSYS,
	"SIGTERM":   syscall.SIGTERM,
	"SIGTRAP":   syscall.SIGTRAP,
	"SIGTSTP":   syscall.SIGTSTP,
	"SIGTTIN":   syscall.SIGTTIN,
	"SIGTTOU":   syscall.SIGTTOU,
	"SIGURG":    syscall.SIGURG,
	"SIGUSR1":   syscall.SIGUSR1,
	"SIGUSR2":   syscall.SIGUSR2,
	"SIGVTALRM": syscall.SIGVTALRM,
	"SIGWINCH":  syscall.SIGWINCH,
	"SIGXCPU":   syscall.SIGXCPU,
	"SIGXFSZ":   syscall.SIGXFSZ,

	// Without SIG prefix
	"ABRT":   syscall.SIGABRT,
	"ALRM":   syscall.SIGALRM,
	"BUS":    syscall.SIGBUS,
	"CHLD":   syscall.SIGCHLD,
	"CONT":   syscall.SIGCONT,
	"EMT":    syscall.SIGEMT,
	"FPE":    syscall.SIGFPE,
	"HUP":    syscall.SIGHUP,
	"ILL":    syscall.SIGILL,
	"INFO":   syscall.SIGINFO,
	"INT":    syscall.SIGINT,
	"IO":     syscall.SIGIO,
	"IOT":    syscall.SIGIOT,
	"KILL":   syscall.SIGKILL,
	"PIPE":   syscall.SIGPIPE,
	"PROF":   syscall.SIGPROF,
	"QUIT":   syscall.SIGQUIT,
	"SEGV":   syscall.SIGSEGV,
	"STOP":   syscall.SIGSTOP,
	"SYS":    syscall.SIGSYS,
	"TERM":   syscall.SIGTERM,
	"TRAP":   syscall.SIGTRAP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"URG":    syscall.SIGURG,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"VTALRM": syscall.SIGVTALRM,
	"WINCH":  syscall.SIGWINCH,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,

	// Numeric values
	// SIGIOT removed as it's the same as SIGABRT
	"6":  syscall.SIGABRT,
	"14": syscall.SIGALRM,
	"10": syscall.SIGBUS,
	"20": syscall.SIGCHLD,
	"19": syscall.SIGCONT,
	"7":  syscall.SIGEMT,
	"8":  syscall.SIGFPE,
	"1":  syscall.SIGHUP,
	"4":  syscall.SIGILL,
	"29": syscall.SIGINFO,
	"2":  syscall.SIGINT,
	"23": syscall.SIGIO,
	"9":  syscall.SIGKILL,
	"13": syscall.SIGPIPE,
	"27": syscall.SIGPROF,
	"3":  syscall.SIGQUIT,
	"11": syscall.SIGSEGV,
	"17": syscall.SIGSTOP,
	"12": syscall.SIGSYS,
	"15": syscall.SIGTERM,
	"5":  syscall.SIGTRAP,
	"18": syscall.SIGTSTP,
	"21": syscall.SIGTTIN,
	"22": syscall.SIGTTOU,
	"16": syscall.SIGURG,
	"30": syscall.SIGUSR1,
	"31": syscall.SIGUSR2,
	"26": syscall.SIGVTALRM,
	"28": syscall.SIGWINCH,
	"24": syscall.SIGXCPU,
	"25": syscall.SIGXFSZ,
}

type signalValue syscall.Signal

func (s *signalValue) Set(value string) error {
	v, ok := signalTable[value]
	if !ok {
		return fmt.Errorf("%s is not a valid signal", value)
	}

	*s = (signalValue)(v)

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
