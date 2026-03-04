package main

import (
	"flag"
	"fmt"
	"os"

	"dailynotes/internal/notes"
)

func main() {
	copyTasks := flag.Bool("copy-tasks", false, "Copy incomplete tasks from previous note to clipboard")
	flag.Parse()

	if *copyTasks {
		if err := notes.CopyIncompleteTasksToClipboard(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := notes.CreateDailyNote(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
