# Log Parser

A Go utility for parsing job execution logs and reporting job durations with threshold-based alerts.

## Overview

This tool parses CSV-formatted log files containing job execution records and provides duration analysis with configurable warning and error thresholds. It's designed to help monitor job performance and identify jobs that are taking longer than expected to complete.

## Features

- Parses CSV log files with job start/end timestamps
- Calculates job durations with automatic midnight wrap-around handling
- Provides threshold-based alerting (warnings and errors)
- Identifies incomplete jobs (missing START or END records)
- Supports time-only format (HH:MM:SS) for intraday job monitoring

## Installation

### Prerequisites
- Go 1.16 or later

### Build from source
```bash
go build -o log_parser main.go
```

## Usage

```bash
./log_parser <log_file_path>
```

### Example
```bash
./log_parser /path/to/logfile
```

## Log File Format

The tool expects CSV files with the following format:

```
timestamp,description,status,pid
```

### Fields
- **timestamp**: Time in HH:MM:SS format (24-hour)
- **description**: Job description (free text)
- **status**: Either "START" or "END"
- **pid**: Process ID or unique job identifier

### Example Log File
```csv
09:15:30,Database Backup,START,PID001
09:17:45,Database Backup,END,PID001
10:30:00,Data Processing,START,PID002
10:35:30,Data Processing,END,PID002
11:45:00,Report Generation,START,PID003
12:02:15,Report Generation,END,PID003
```

## Output Format

The tool provides different output based on job duration:

### Normal Jobs
```
Job PID001 duration: 2m15s
```

### Warning Threshold (> 5 minutes)
```
Warning: Job PID002 took longer than 5 minutes: 7m30s
```

### Error Threshold (> 10 minutes)
```
Error: Job PID003 took longer than 10 minutes: 17m15s
```

### Incomplete Jobs
```
Incomplete job PID004: missing END
Incomplete job PID005: missing START
```

## Special Features

### Midnight Wrap-around
The tool automatically handles jobs that span midnight by adding 24 hours to the duration calculation when the end time appears to be before the start time.

Example:
- Start: 23:45:00
- End: 00:15:00
- Calculated Duration: 30m0s (not -23h30m)

### Error Handling
- Invalid log lines are reported but don't stop processing
- Malformed timestamps are logged with line numbers
- Unknown status values are flagged
- File access errors are handled gracefully

## Thresholds

The tool uses fixed thresholds for job duration analysis:

- **Normal**: ≤ 5 minutes
- **Warning**: > 5 minutes and ≤ 10 minutes  
- **Error**: > 10 minutes

## Exit Codes

- `0`: Success
- `1`: Invalid usage or file access error

## Limitations

- Only supports time-only format (HH:MM:SS), not full datetime
- Assumes all jobs complete within a 24-hour period
- Thresholds are hardcoded (not configurable via command line)
- CSV parsing is basic (doesn't handle quoted fields with commas)

## Example Use Cases

- Monitor cron job execution times
- Analyze batch processing performance
- Identify jobs that may be hanging or performing poorly
- Generate daily/weekly job performance reports

## Troubleshooting

### Common Issues

1. **"Invalid log line" messages**: Check that your CSV has exactly 4 columns
2. **"Error parsing timestamp" messages**: Ensure timestamps are in HH:MM:SS format
3. **"Incomplete job" messages**: Verify that each job has both START and END entries with matching PIDs

### Debug Tips

- Use line numbers in error messages to locate problematic entries
- Verify CSV format matches expected structure
- Check for extra whitespace or special characters in PID fields

## To Improve:
 - Report generator as JSON file
 - Multiple files ingested at once
 - Email notification based on LogLevel 
