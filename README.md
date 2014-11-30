## MultiBar

Display progress bars in Go

    $ go run main.go
    here we have a progress bar
    some work  30% [====================>-------------------------------------------] 925ms
    and here we have another progress bar
    here we have a longer prepend string  25% [============>-----------------------------------] 911ms
    and a third
    working...  19% [=============>--------------------------------------------------] 911ms

Example usage:

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

        // some arbitrary totals for our three progress bars
        // in practice, these could be file sizes or similar
        first, second, third := 150, 100, 200

        // make some output for the screen
        // the MakeBar(total, prependString) returns a method that you can pass progress into
        fmt.Println("here we have a progress bar")
        barProgress1 := progressBars.MakeBar(first, "some work")

        fmt.Println("and here we have another progress bar")
        barProgress2 := progressBars.MakeBar(second, "here we have a longer prepend string")

        fmt.Println("and a third")
        barProgress3 := progressBars.MakeBar(third, "working...")

        // spawn the listener that will allow the progres bars to update
        // I should be able to move this into the New() method
        go progressBars.Listen()

        // mock work. spawn three goroutines to do arbitrary work, updating their
        // respective progress bars as they see fit
        wg := &sync.WaitGroup{}
        wg.Add(3)
        go func() {
            // do something asyn that we can get updates upon
            // every time an update comes in, tell the bar to re-draw
            for i := 0; i <= first; i++ {
                barProgress1(i)
                time.Sleep(time.Millisecond * 15)
            }
            wg.Done()
        }()

        go func() {
            for i := 0; i <= second; i++ {
                barProgress2(i)
                time.Sleep(time.Millisecond * 25)
            }
            wg.Done()
        }()

        go func() {
            for i := 0; i <= third; i++ {
                barProgress3(i)
                time.Sleep(time.Millisecond * 12)
            }
            wg.Done()
        }()

        wg.Wait()

        // continue doing other work
        fmt.Println("All Bars Complete")
    }
