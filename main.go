package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type JobStatus struct {
	Start string
	End   string
}

type JobEntry struct {
	Timestamp   string
	Description string
	Status      string
	PID         string
}

const timeLayout = "15:04:05"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: log_parser <log_file>")
		os.Exit(1)
	}
	logFilePath := os.Args[1]
	file, err := os.Open(logFilePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
}

func parseLine(line string) []string {
	return strings.Split(line, ",")
}

func calculateDuration(start, end string) (time.Duration, error) {
	if start == "" || end == "" {
		return 0, fmt.Errorf("start or end time is empty")
	}
	startTime, err := time.Parse(timeLayout, start)
	if err != nil {
		return 0, fmt.Errorf("error parsing start time: %v", err)
	}
	endTime, err := time.Parse(timeLayout, end)
	if err != nil {
		return 0, fmt.Errorf("error parsing end time: %v", err)
	}
	duration := endTime.Sub(startTime)
	if duration < 0 {
		duration += 24 * time.Hour
	}
	return duration, nil
}
