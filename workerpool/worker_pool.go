package workerpool

import (
	"fmt"
	"sync"

	"rincon-orlando/go-bootcamp/config"
)

type GoRoutinePool struct {
	queue  chan work
	wg     sync.WaitGroup
	config config.GoRoutinePoolConfig
}

type workFunc interface {
	Run(ch chan<- interface{}) bool
}

type work struct {
	fn workFunc
}

func New(numWorkers int, config config.GoRoutinePoolConfig) *GoRoutinePool {
	gp := &GoRoutinePool{
		// Had to make this a buffered channel, otherwise this was blocking on some requests like
		// http://localhost:8082/pokemons/filter?type=even&items=30&items_per_workers=1
		// or
		// http://localhost:8082/pokemons/filter?type=even&items=6&items_per_workers=1
		queue:  make(chan work, 20),
		config: config,
	}

	gp.AddWorkers(numWorkers)
	return gp
}

func (gp *GoRoutinePool) ScheduleWork(fn workFunc) {
	gp.queue <- work{fn}
}

func (gp *GoRoutinePool) AddWorkers(numWorkers int) {
	gp.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			items_per_worker := gp.config.ItemsPerWorker
			for job := range gp.queue {
				if job.fn.Run(gp.config.Channel) {
					items_per_worker--
					// This worker has finished its search because it was capped to certain items per worker
					if items_per_worker == 0 {
						// fmt.Printf("Worker %d found all its corresponding items\n", workerID)
						break
					}
				}
			}
			fmt.Printf("Worker %d is done\n", workerID)
			gp.wg.Done()                        // Complete this wait group task
			gp.config.DoneChannel <- struct{}{} // Send a done signal for this
		}(i)
	}
}

func (gp *GoRoutinePool) Monitor() []interface{} {
	response := make([]interface{}, 0)
	foundItems := 0

	for n := gp.config.NumWorkers; n > 0; {
		select {
		case entry := <-gp.config.Channel:
			fmt.Printf("New item arrived %s\n", entry)
			response = append(response, entry)
			foundItems++
			// Case: reached the amount of valid items you need to display as a response
			if foundItems == gp.config.Items {
				// Look no more, we found the amount of pokemons requested
				fmt.Printf("Desired number of items [%d] obtained\n", gp.config.Items)
				n = 0 // Set this to 0 to break the outter loop
			}
		case <-gp.config.DoneChannel:
			n--
		}
	}

	return response
}

func (gp *GoRoutinePool) Close() {
	close(gp.queue)
	gp.wg.Wait()
}
