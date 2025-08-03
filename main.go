package main

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
