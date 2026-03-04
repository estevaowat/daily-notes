package main

import (
	"fmt"
	"os"

	"dailynotes/internal/notes"
)

func main() {
	if err := notes.CreateDailyNote(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
