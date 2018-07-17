package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Returns a context who's cancel() method is called when the given signals are received
func SigContext() context.Context {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, // https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
	)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()
	return ctx
}
