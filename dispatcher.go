package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
)

// WorkerQueue does stuff
var WorkerQueue chan chan WorkRequest

// StartDispatcher does stuff
func StartDispatcher(nworkers int, workspace string) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	history := make([]int, 0)

	files, err := ioutil.ReadDir(workspace)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println("Found file: " + f.Name())

		if f.IsDir() != true {
			fmt.Println("Skipping file: " + f.Name())
			continue
		}

		h, err := strconv.Atoi(f.Name())

		if err != nil {
			fmt.Println("DEBUG: Could not index " + f.Name())
			continue
		}

		fmt.Printf("DEBUG: Adding " + f.Name() + " to history array\n")

		history = append(history, h)
	}

	sort.Ints(history)

	fmt.Println("Sorted history: ", history)

	startpoint := history[len(history)-1]

	startpoint++

	fmt.Println("Starting with: ", startpoint)

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
				fmt.Println("Received work requeust from webserver...", startpoint)
				go func() {
					worker := <-WorkerQueue

					fmt.Println("Dispatching work request to worker...")
					worker <- work

				}()
			}
		}
	}()
}
