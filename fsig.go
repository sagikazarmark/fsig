package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/fsnotify/fsnotify"
)

// Provisioned by ldflags
var (
	Version string
)

var (
	watch = kingpin.Flag("watch", "Watched directory (at least one)").Short('w').Required().Strings()
	sig   = signalArg(kingpin.Arg("signal", "Signal to be sent to the child process").Required())
	cmd   = kingpin.Arg("cmd", "Child process command").Required().String()
	args  = kingpin.Arg("args", "Child process arguments").Strings()
)

func init() {
	kingpin.CommandLine.Name = "fsig"
	kingpin.CommandLine.Help = "Send signals to a child process upon file changes"
	kingpin.CommandLine.Version(Version)
}

func main() {
	kingpin.Parse()

	watcher := newWatcher(*watch)
	childCmd := newChildCommand(*cmd, *args)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	err := childCmd.Start()
	if err != nil {
		log.Fatalln("error:", err)
	}

	go func() {
		for {
			select {

			case <-done: // Process exited
				return

			case sig := <-signals: // Forwarding signal to child process
				log.Println("received signal:", sig)

				err := childCmd.Process.Signal(sig)
				if err != nil {
					fail(childCmd, err)
				}

			case event := <-watcher.Events: // Change detected
				if event.Op&fsnotify.Create == fsnotify.Create { // TODO: which changes should be watched?
					log.Println("received event:", event)

					err := childCmd.Process.Signal(*sig)
					if err != nil {
						fail(childCmd, err)
					}
				}

			case err := <-watcher.Errors: // Error watching changes
				fail(childCmd, err)
			}
		}
	}()

	err = childCmd.Wait()
	close(done)

	if err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			os.Exit(eerr.Sys().(syscall.WaitStatus).ExitStatus())
		}

		log.Fatalln("error:", err)
	}
}

func newWatcher(watch []string) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatalln(err)
	}

	for _, w := range watch {
		err := watcher.Add(w)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return watcher
}

func newChildCommand(cmd string, args []string) *exec.Cmd {
	childCmd := exec.Command(cmd, args...)

	childCmd.Env = os.Environ()
	childCmd.Stdin = os.Stdin
	childCmd.Stdout = os.Stdout
	childCmd.Stderr = os.Stderr

	return childCmd
}

func fail(cmd *exec.Cmd, err error) {
	e := cmd.Process.Kill()
	if e != nil {
		log.Println("process exit error:", err)
	}

	log.Fatalln(err)
}
