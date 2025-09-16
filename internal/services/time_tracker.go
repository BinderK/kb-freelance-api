package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"kb-freelance-api/internal/config"
)

type TimeTrackerService struct {
	config *config.Config
}

func NewTimeTrackerService(cfg *config.Config) *TimeTrackerService {
	return &TimeTrackerService{config: cfg}
}

type TimeEntry struct {
	ID              int        `json:"id"`
	Client          string     `json:"client"`
	Project         string     `json:"project"`
	Description     string     `json:"description"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	DurationMinutes int        `json:"duration_minutes"`
	IsRunning       bool       `json:"is_running"`
}

type TimerStatus struct {
	IsRunning       bool      `json:"is_running"`
	Client          string    `json:"client,omitempty"`
	Project         string    `json:"project,omitempty"`
	Description     string    `json:"description,omitempty"`
	StartTime       time.Time `json:"start_time,omitempty"`
	DurationMinutes int       `json:"duration_minutes,omitempty"`
}

type TodaySummary struct {
	TotalHours   float64     `json:"total_hours"`
	TotalMinutes int         `json:"total_minutes"`
	EntryCount   int         `json:"entry_count"`
	Breakdown    []Breakdown `json:"breakdown"`
	RawOutput    string      `json:"raw_output,omitempty"`
}

type Breakdown struct {
	ClientProject string  `json:"client_project"`
	Hours         float64 `json:"hours"`
	Minutes       int     `json:"minutes"`
}

func (s *TimeTrackerService) StartTimer(client, project, description string) (map[string]interface{}, error) {
	// Build command to start timer using configurable Python executable
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "start", client, project)
	if description != "" {
		cmd.Args = append(cmd.Args, "--desc", description)
	}

	// Set working directory to time tracker path
	cmd.Dir = s.config.TimeTrackerPath

	// Debug: log the command and directory
	fmt.Printf("DEBUG: Running command: %v\n", cmd.Args)
	fmt.Printf("DEBUG: Working directory: %s\n", cmd.Dir)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to start timer: %s, output: %s", err.Error(), string(output))
	}

	// Return the timer data that was just started
	now := time.Now()
	return map[string]interface{}{
		"id":               1, // This should be the actual ID from the database
		"client":           client,
		"project":          project,
		"description":      description,
		"start_time":       now.Format(time.RFC3339),
		"is_running":       true,
		"duration_minutes": 0,
	}, nil
}

func (s *TimeTrackerService) StopTimer() (map[string]interface{}, error) {
	// Build command to stop timer using configurable Python executable
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "stop")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to stop timer: %s, output: %s", err.Error(), string(output))
	}

	// After stopping the timer, return the stopped timer data
	// For now, return a proper TimeEntry structure indicating it's stopped
	now := time.Now()
	return map[string]interface{}{
		"id":               1,                                        // This should be the actual ID from the database
		"client":           "Unknown",                                // This should be the actual client from the database
		"project":          "Unknown",                                // This should be the actual project from the database
		"description":      "",                                       // This should be the actual description from the database
		"start_time":       now.Add(-time.Hour).Format(time.RFC3339), // This should be the actual start time
		"end_time":         now.Format(time.RFC3339),
		"is_running":       false,
		"duration_minutes": 60, // This should be the actual duration
	}, nil
}

func (s *TimeTrackerService) GetStatus() (map[string]interface{}, error) {
	// Build command to get status using configurable Python executable with JSON output
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "status", "--json")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If no timer is running, this is not an error
		if len(output) > 0 {
			outputStr := string(output)
			if contains(outputStr, "null") {
				return nil, nil // Return nil to indicate no timer running
			}
		}
		return nil, fmt.Errorf("failed to get timer status: %s, output: %s", err.Error(), string(output))
	}

	// Parse JSON output
	outputStr := string(output)
	fmt.Printf("DEBUG: GetStatus output: %s\n", outputStr)

	if outputStr == "null" || outputStr == "" {
		return nil, nil // No timer running
	}

	// Parse the JSON response
	var timerData map[string]interface{}
	if err := json.Unmarshal(output, &timerData); err != nil {
		return nil, fmt.Errorf("failed to parse timer status JSON: %s, output: %s", err.Error(), string(output))
	}

	fmt.Printf("DEBUG: Parsed timer data: %+v\n", timerData)
	return timerData, nil
}

func (s *TimeTrackerService) GetRecentEntries(limit int) ([]TimeEntry, error) {
	// Build command to get recent entries using configurable Python executable with JSON output
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "list", "--json")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent entries: %s, output: %s", err.Error(), string(output))
	}

	// Parse JSON output
	fmt.Printf("DEBUG: GetRecentEntries output: %s\n", string(output))
	var entriesData []map[string]interface{}
	if err := json.Unmarshal(output, &entriesData); err != nil {
		return nil, fmt.Errorf("failed to parse recent entries JSON: %s, output: %s", err.Error(), string(output))
	}
	fmt.Printf("DEBUG: Parsed entries data: %+v\n", entriesData)

	// Convert to TimeEntry structs
	var entries []TimeEntry
	for _, item := range entriesData {
		entry := TimeEntry{
			ID:              int(item["id"].(float64)),
			Client:          item["client"].(string),
			Project:         item["project"].(string),
			Description:     item["description"].(string),
			DurationMinutes: int(item["duration_minutes"].(float64)),
			IsRunning:       item["is_running"].(bool),
		}

		// Parse start time
		if startTimeStr, ok := item["start_time"].(string); ok {
			// Try different time formats
			if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
				entry.StartTime = startTime
			} else if startTime, err := time.Parse(time.RFC3339Nano, startTimeStr); err == nil {
				entry.StartTime = startTime
			} else if startTime, err := time.Parse("2006-01-02T15:04:05.999999", startTimeStr); err == nil {
				entry.StartTime = startTime
			} else {
				fmt.Printf("DEBUG: Failed to parse start_time: %s, error: %v\n", startTimeStr, err)
			}
		}

		// Parse end time if it exists
		if endTimeStr, ok := item["end_time"].(string); ok && endTimeStr != "" {
			// Try different time formats
			if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
				entry.EndTime = &endTime
			} else if endTime, err := time.Parse(time.RFC3339Nano, endTimeStr); err == nil {
				entry.EndTime = &endTime
			} else if endTime, err := time.Parse("2006-01-02T15:04:05.999999", endTimeStr); err == nil {
				entry.EndTime = &endTime
			} else {
				fmt.Printf("DEBUG: Failed to parse end_time: %s, error: %v\n", endTimeStr, err)
			}
		}

		entries = append(entries, entry)
	}

	// Apply limit if specified
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}

	return entries, nil
}

func (s *TimeTrackerService) GetTodaySummary() (*TodaySummary, error) {
	// Build command to get today's summary using configurable Python executable with JSON output
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "today", "--json")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get today's summary: %s, output: %s", err.Error(), string(output))
	}

	// Parse JSON output
	var summaryData map[string]interface{}
	if err := json.Unmarshal(output, &summaryData); err != nil {
		return nil, fmt.Errorf("failed to parse today's summary JSON: %s, output: %s", err.Error(), string(output))
	}

	// Extract data from JSON
	totalHours, _ := summaryData["total_hours"].(float64)
	totalMinutes, _ := summaryData["total_minutes"].(float64)
	entryCount, _ := summaryData["entry_count"].(float64)
	breakdownData, _ := summaryData["breakdown"].([]interface{})

	// Convert breakdown data
	var breakdown []Breakdown
	for _, item := range breakdownData {
		if breakdownItem, ok := item.(map[string]interface{}); ok {
			client, _ := breakdownItem["client"].(string)
			project, _ := breakdownItem["project"].(string)
			durationMinutes, _ := breakdownItem["duration_minutes"].(float64)

			breakdown = append(breakdown, Breakdown{
				ClientProject: fmt.Sprintf("%s - %s", client, project),
				Hours:         durationMinutes / 60,
				Minutes:       int(durationMinutes),
			})
		}
	}

	return &TodaySummary{
		TotalHours:   totalHours,
		TotalMinutes: int(totalMinutes),
		EntryCount:   int(entryCount),
		Breakdown:    breakdown,
		RawOutput:    string(output),
	}, nil
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
