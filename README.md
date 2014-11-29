## MultiBar

Display progress bars in Go

    $ go run main.go
    Inline Progress Bar:   7% [=====>------------------------------------------------------] 309ms

Example usage:

package main

    import (
        "fmt"
        "time"

        "github.com/sethgrid/multibar"
    )

    func main() {
        total := 150

        bar, _ := multibar.New(total)
        bar.Prepend("Inline Progress Bar: ")

        // display a progress bar
        for i := 0; i <= total; i++ {
            bar.Display(i)
            time.Sleep(time.Millisecond * 25)
        }
        // end the previous last line of output
        fmt.Println()
        fmt.Println("Complete")
    }