package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
)

// FindHistory finds the first point
func FindHistory(workspace string) int {

	history := make([]int, 0)

	files, err := ioutil.ReadDir(workspace)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println("Found file: " + f.Name())

		if f.IsDir() != true {
			fmt.Println("DEBUG: Skipping file - " + f.Name())
			continue
		}

		h, err := strconv.Atoi(f.Name())

		if err != nil {
			fmt.Println("DEBUG: Could not index - " + f.Name())
			continue
		}

		fmt.Printf("DEBUG: Adding - " + f.Name() + " to history array\n")

		history = append(history, h)
	}

	sort.Ints(history)

	fmt.Println("DEBUG: Sorted history - ", history)

	startpoint := history[len(history)-1]

	fmt.Println("DEBUG: History startpoint: ", startpoint)

	return startpoint
}
