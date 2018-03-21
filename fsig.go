package main

import (
	"fmt"
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
	child := newChild(*cmd, *args)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL) // TODO: which signals should be forwarded?

	done := make(chan struct{})

	err := child.Start()
	if err != nil {
		log.Fatalln("error:", err)
	}

	go func() {
		for {
			select {

			case <-done: // Process exited
				return

			case sig := <-signals: // Forwarding signal to child process
				err := child.Process.Signal(sig)
				if err != nil {
					log.Println("error:", err)
					child.Process.Kill()
				}

			case event := <-watcher.Events: // Change detected
				if event.Op&fsnotify.Create == fsnotify.Create { // TODO: which changes should be watched?
					err := child.Process.Signal(*sig)
					if err != nil {
						log.Println("error:", err)
						child.Process.Kill()
					}
				}

			case err := <-watcher.Errors: // Error watching changes
				log.Println("error:", err)
				child.Process.Kill()
			}
		}
	}()

	err = child.Wait()
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
		fmt.Println(err)

		os.Exit(2)
	}

	for _, w := range watch {
		watcher.Add(w)
	}

	return watcher
}

func newChild(cmd string, args []string) *exec.Cmd {
	childCmd := exec.Command(cmd, args...)

	childCmd.Env = os.Environ()
	childCmd.Stdin = os.Stdin
	childCmd.Stdout = os.Stdout
	childCmd.Stderr = os.Stderr

	return childCmd
}
