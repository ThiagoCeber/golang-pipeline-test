package stage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
)

// First pipeline stage.
func First(ctx context.Context, quit chan struct{}, done *sync.WaitGroup) <-chan string {
	output := make(chan string)
	done.Add(1)
	var count int

	go func() {
		defer done.Done()
		for {
			select {
			case <-quit:
				close(output)
				fmt.Println("stage 1 closed")
				return
			default:
				count++
				output <- strconv.Itoa(count)
				log.Println(fmt.Sprintf("message %v created.", count))
			}
		}
	}()

	return output
}
