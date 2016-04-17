package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sethgrid/multibar"
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		log.Fatal("Exit with incomplete download (expected behavior)")
		os.Exit(1)
	}()

	progressBars, _ := multibar.New()
	progressBars.Println("Below is a progress bar.")

	sleepTime := 100 * time.Millisecond
	iterations := 100
	totalTimeToComplete := time.Duration(iterations) * sleepTime
	timeUnilSignal := totalTimeToComplete / time.Duration(2)
	barProgress1 := progressBars.MakeBar(iterations, "FooBar:")

	go progressBars.Listen()

	wg := &sync.WaitGroup{}
	wg.Add(iterations)
	go func() {
		// do something asyn that we can get updates upon
		// every time an update comes in, tell the bar to re-draw
		// this could be based on transferred bytes or similar
		for i := 0; i <= iterations; i++ {
			barProgress1(i)
			time.Sleep(time.Millisecond * 100)
			wg.Done()
		}
	}()

	<-time.After(timeUnilSignal)

	err := syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	if err != nil {
		log.Println("got an error trying to kill self with SIGHUP", err)
	}

	wg.Wait()
	log.Println("program completed. That should not have happened.")
}
