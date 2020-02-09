package main

import (
	"fmt"
)

// WorkerQueue does stuff
var WorkerQueue chan chan WorkRequest

// StartDispatcher does stuff
func StartDispatcher(nworkers int, workspace string) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, workspace, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work requeust from webserver...")
				go func() {
					worker := <-WorkerQueue

					fmt.Println("Dispatching work request to worker...")
					worker <- work

				}()
			}
		}
	}()
}
