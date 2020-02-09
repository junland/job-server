package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// WorkQueue A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

// WorkRequest does stuff
type WorkRequest struct {
	WorkID  int
	WorkDir string
	Name    string
	Delay   time.Duration
}

// Collector jkh
func (c *Config) Collector(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the delay.
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check to make sure the delay is anywhere from 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Now, we retrieve the person's name from the request.
	name := r.FormValue("name")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	c.CurrentHistory++

	fmt.Println("Added workid of ", c.CurrentHistory)

	id := strconv.Itoa(c.CurrentHistory)

	// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
	work := WorkRequest{WorkID: c.CurrentHistory, Name: name, Delay: delay, WorkDir: c.WorkSpaceDir + "/" + id}

	// Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued - ")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}

// StopWorker jkh
func (c *Config) StopWorker(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Now, we retrieve the person's name from the request.
	workername := r.FormValue("worker")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if workername == "" {
		http.Error(w, "You must specify a worker to stop.", http.StatusBadRequest)
		return
	}

	fmt.Println("Work request queued - ")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}

// curl localhost:8000/work -d name=john -d delay=60s && curl localhost:8000/stop -d worker=worker1
