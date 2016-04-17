package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/sethgrid/multibar"
)

func main() {
	// create the multibar container
	// this allows our bars to work together without stomping on one another
	progressBars, _ := multibar.New()

	// some arbitrary totals for our  progress bars
	// in practice, these could be file sizes or similar
	mediumTotal, smallTotal, largerTotal := 150, 100, 200

	// make some output for the screen
	// the MakeBar(total, prependString) returns a method that you can pass progress into
	progressBars.Println("Below are many progress bars.")
	progressBars.Println("It is best to use the print wrappers to keep output synced up.")
	progressBars.Println("We can switch back to normal fmt after our progress bars are done.\n")

	// we will update the progress down below in the mock work section with barProgress1(int)
	barProgress1 := progressBars.MakeBar(mediumTotal, "1st")

	progressBars.Println()
	progressBars.Println("We can separate bars with blocks of text, or have them grouped.\n")

	barProgress2 := progressBars.MakeBar(smallTotal, "2nd - with description:")
	barProgress3 := progressBars.MakeBar(largerTotal, "3rd")
	barProgress4 := progressBars.MakeBar(mediumTotal, "4th")
	barProgress5 := progressBars.MakeBar(smallTotal, "5th")
	barProgress6 := progressBars.MakeBar(largerTotal, "6th")

	progressBars.Println("And we can have blocks of text as we wait for progress bars to complete...")

	// listen in for changes on the progress bars
	// I should be able to move this into the constructor at some point
	go progressBars.Listen()

	/*

	   *** mock work ***
	   spawn some goroutines to do arbitrary work, updating their
	   respective progress bars as they see fit

	*/
	wg := &sync.WaitGroup{}
	wg.Add(6)
	go func() {
		// do something asyn that we can get updates upon
		// every time an update comes in, tell the bar to re-draw
		// this could be based on transferred bytes or similar
		for i := 0; i <= mediumTotal; i++ {
			barProgress1(i)
			time.Sleep(time.Millisecond * 15)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i <= smallTotal; i++ {
			barProgress2(i)
			time.Sleep(time.Millisecond * 25)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i <= largerTotal; i++ {
			barProgress3(i)
			time.Sleep(time.Millisecond * 12)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress4(i)
			time.Sleep(time.Millisecond * 10)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= smallTotal; i++ {
			barProgress5(i)
			time.Sleep(time.Millisecond * 20)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= largerTotal; i++ {
			barProgress6(i)
			time.Sleep(time.Millisecond * 10)
		}
		wg.Done()
	}()
	wg.Wait()

	// continue doing other work
	fmt.Println("All Bars Complete")
}
