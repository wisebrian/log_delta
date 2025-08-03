package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type JobStatus struct {
	Start time.Time
	End   time.Time
}

type JobEntry struct {
	Timestamp   time.Time
	Description string
	Status      string
	PID         string
}

const timeLayout = "15:04:05"

func main() {
	// Expect exactly one command-line argument: the path to the log file
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
	// Map to store JobStatus keyed by PID
	entries := make(map[string]*JobStatus)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		parts := parseLine(line)
		// Ensure there are exactly 4 components: Timestamp, Description, Status, PID
		if len(parts) != 4 {
			fmt.Printf("Invalid log line: %s\n", line)
			continue
		}
		// Parse timestamp to standardize the format
		rawTime := strings.TrimSpace(parts[0])
		parsedTime, err := time.Parse(timeLayout, rawTime)
		if err != nil {
			fmt.Printf("Line %d: Error parsing timestamp '%s': %v\n", lineNumber, rawTime, err)
			continue
		}

		// Construct a JobEntry struct from line parts
		job := JobEntry{
			Timestamp:   parsedTime,
			Description: strings.TrimSpace(parts[1]),
			Status:      strings.TrimSpace(parts[2]),
			PID:         strings.TrimSpace(parts[3]),
		}

		// Ensure the job entry exists in the map
		if _, ok := entries[job.PID]; !ok {
			entries[job.PID] = &JobStatus{}
		}
		switch job.Status {
		case "START":
			entries[job.PID].Start = job.Timestamp
		case "END":
			entries[job.PID].End = job.Timestamp
		default:
			fmt.Printf("Line %d: Status Unknown '%s'\n", lineNumber, job.Status)
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	// Report of all jobs
	for pid, status := range entries {
		reportJobDuration(pid, status)
	}
}

// parseLine splits a CSV log line into its parts
func parseLine(line string) []string {
	return strings.Split(line, ",")
}

// calculateDuration diffs the time difference between start and end times
// Handles wrap-around at midnight by adding 24 hours if needed
func calculateDuration(start, end time.Time) time.Duration {
	/*if start == "" || end == "" {
		return 0, fmt.Errorf("start or end time is empty")
	}
	startTime, err := time.Parse(timeLayout, start)
	if err != nil {
		return 0, fmt.Errorf("error parsing start time: %v", err)
	}
	endTime, err := time.Parse(timeLayout, end)
	if err != nil {
		return 0, fmt.Errorf("error parsing end time: %v", err)
	}*/
	duration := end.Sub(start)
	if duration < 0 {
		duration += 24 * time.Hour
	}
	return duration
}

func reportJobDuration(pid string, status *JobStatus) {
	var missingField string
	if status.Start.IsZero() {
		missingField = "START"
	} else if status.End.IsZero() {
		missingField = "END"
	}
	if missingField != "" {
		fmt.Printf("Incomplete job %s: missing %s\n", pid, missingField)
		return
	}
	duration := calculateDuration(status.Start, status.End)

	switch {
	case duration > 10*time.Minute:
		fmt.Printf("Error: Job %s took longer than 10 minutes: %s\n", pid, duration)
	case duration > 10*time.Minute:
		fmt.Printf("Warning: Job %s took longer than 5 minutes: %s\n", pid, duration)
	default:
		fmt.Printf("Job %s duration: %s\n", pid, duration)
	}
}
