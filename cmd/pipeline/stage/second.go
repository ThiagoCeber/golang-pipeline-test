package stage

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// Second pipeline stage.
func Second(ctx context.Context, maxWorkers int, input <-chan string, done *sync.WaitGroup) <-chan string {
	output := make(chan string)
	sem := semaphore.NewWeighted(int64(maxWorkers))
	done.Add(1)

	go func() {
		defer func() {
			log.Println("stage 2 done")
			done.Done()
		}()
		for message := range input {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Println(fmt.Errorf("failed to acquire semaphore resource. [%v]", err))
				break
			}
			go func(message string) {
				defer sem.Release(1)

				start := time.Now()
				time.Sleep(5 * time.Second)

				output <- fmt.Sprintf("%s --> [stage2(%v)]", message, time.Since(start).Seconds())

			}(message)
		}
		if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
			log.Println(fmt.Errorf("failed to acquire all semaphore resources. [%v]", err))
		}
		close(output)
		log.Println("stage 2 closed")
	}()

	return output
}
