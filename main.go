package main

import (
	"bufio"
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
	entries := make(map[string]*JobStatus)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := parseLine(line)
		if len(parts) != 4 {
			fmt.Printf("Invalid log line: %s\n", line)
			continue
		}
		job := JobEntry{
			Timestamp:   strings.TrimSpace(parts[0]),
			Description: strings.TrimSpace(parts[1]),
			Status:      strings.TrimSpace(parts[2]),
			PID:         strings.TrimSpace(parts[3]),
		}

		parsedTime, err := time.Parse(timeLayout, job.Timestamp)
		if err != nil {
			fmt.Printf("Error parsing timestamp: %v\n", err)
			continue
		}
		job.Timestamp = parsedTime.Format(timeLayout)
		if _, ok := entries[job.PID]; !ok {
			entries[job.PID] = &JobStatus{}
		}
		if job.Status == "START" {
			entries[job.PID].Start = job.Timestamp
		}
		if job.Status == "END" {
			entries[job.PID].End = job.Timestamp
		}
		if entries[job.PID].Start != "" && entries[job.PID].End != "" {
			d, err := calculateDuration(entries[job.PID].Start, entries[job.PID].End)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if d > 5*time.Minute && d < 10*time.Minute {
				fmt.Printf("Warning: Job %s took longer than 5 minutes: %s\n", job.PID, d)
			} else if d > 10*time.Minute {
				fmt.Printf("Error: Job %s took longer than 10 minutes: %s\n", job.PID, d)
			}
		}
	}
	for pid, status := range entries {
		if status.Start == "" || status.End == "" {
			fmt.Printf("Incomplete job %s: missing %s\n", pid,
				map[bool]string{true: "END", false: "START"}[status.Start != ""])
			continue
		}
		d, err := calculateDuration(status.Start, status.End)
		if err != nil {
			fmt.Printf("Error calculating duration for job %s: %v\n", pid, err)
			continue
		}

		if d > 10*time.Minute {
			fmt.Printf("Error: Job %s took longer than 10 minutes: %s\n", pid, d)
		} else if d > 5*time.Minute {
			fmt.Printf("Warning: Job %s took longer than 5 minutes: %s\n", pid, d)
		} else {
			fmt.Printf("Job %s duration: %s\n", pid, d)
		}
	}
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
