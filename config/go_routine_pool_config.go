package config

import (
	"rincon-orlando/go-bootcamp/util/enum"
)

// GoRoutinePoolConfig - Holds all the relevant information to make the WorkerPool to work
type GoRoutinePoolConfig struct {
	OddEven        enum.OddEven     // Whether this applies for odd or even
	NumWorkers     int              // Num of total workers for this routine pool
	Items          int              // Num of (max) items expected at the end of the call
	ItemsPerWorker int              // (Max) items per worker
	Channel        chan interface{} // Channel where the output information will be written
	DoneChannel    chan struct{}    // Signal channel to indicate finished workers
}

// NewPoolConfig - Factory method
func NewPoolConfig(oe enum.OddEven, numWorkers int, items int, ipw int) GoRoutinePoolConfig {
	// Create a channel so the workers inject elements as they find them
	ch := make(chan interface{}, 50)
	// Create a channel to indicate a worker is done
	doneCh := make(chan struct{}, 10)

	return GoRoutinePoolConfig{oe, numWorkers, items, ipw, ch, doneCh}
}
