package main

import (
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	line := "09:15:30,Database Backup,START,PID001"
	result := parseLine(line)

	if len(result) != 4 {
		t.Errorf("Expected 4 parts, got %d", len(result))
	}

	if result[0] != "09:15:30" {
		t.Errorf("Expected '09:15:30', got '%s'", result[0])
	}
}

func TestCalculateDuration(t *testing.T) {
	start, _ := time.Parse("15:04:05", "09:15:30")
	end, _ := time.Parse("15:04:05", "09:17:45")

	duration := calculateDuration(start, end)
	expected := 2*time.Minute + 15*time.Second

	if duration != expected {
		t.Errorf("Expected %v, got %v", expected, duration)
	}
}

func TestCalculateDurationMidnight(t *testing.T) {
	start, _ := time.Parse("15:04:05", "23:45:00")
	end, _ := time.Parse("15:04:05", "01:15:00")

	duration := calculateDuration(start, end)
	expected := 1*time.Hour + 30*time.Minute

	if duration != expected {
		t.Errorf("Expected %v, got %v", expected, duration)
	}
}
