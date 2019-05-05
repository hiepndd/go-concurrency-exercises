//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"log"
	"os"
	"os/signal"
)

var signals chan os.Signal
var counter int

func main() {

	signals = make(chan os.Signal)

	// Create a process

	proc := MockProcess{}

	// Run the process (blocking)
	go proc.Run()

	signal.Notify(signals, os.Interrupt)

	sig := <-signals

	log.Println("\nSignal received.", sig)

	done := make(chan struct{})

	go proc.Stop(done)

	select {
	case sig = <-signals:
		log.Println("\nSignal received.", sig)
		os.Exit(0)
	case <-done:
		log.Println("\nPorcess shutdown succesfully")
	}

}
