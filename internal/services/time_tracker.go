package services

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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

	return map[string]interface{}{
		"status":      "success",
		"message":     fmt.Sprintf("Started tracking time for %s/%s", client, project),
		"client":      client,
		"project":     project,
		"description": description,
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

	return map[string]interface{}{
		"status":  "success",
		"message": "Timer stopped successfully",
		"output":  string(output),
	}, nil
}

func (s *TimeTrackerService) GetStatus() (*TimerStatus, error) {
	// Build command to get status using configurable Python executable
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "status")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If no timer is running, this is not an error
		if len(output) > 0 {
			outputStr := string(output)
			if contains(outputStr, "No timer is currently running") {
				return &TimerStatus{IsRunning: false}, nil
			}
		}
		return nil, fmt.Errorf("failed to get timer status: %s, output: %s", err.Error(), string(output))
	}

	// Parse the output to extract timer information
	// This is a simplified parser - in a real implementation, you might want to add JSON output to the Python CLI
	outputStr := string(output)

	// For now, we'll return a basic status
	// In a real implementation, you'd parse the actual output
	if contains(outputStr, "Timer Running") {
		return &TimerStatus{
			IsRunning: true,
			// You would parse the actual values from the output here
		}, nil
	}

	return &TimerStatus{IsRunning: false}, nil
}

func (s *TimeTrackerService) GetRecentEntries(limit int) ([]TimeEntry, error) {
	// Build command to get recent entries using configurable Python executable
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "list")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent entries: %s, output: %s", err.Error(), string(output))
	}

	// For now, return empty slice since we need to parse the table output
	// In a real implementation, you'd parse the table output or add JSON output to the Python CLI
	return []TimeEntry{}, nil
}

func (s *TimeTrackerService) GetTodaySummary() (*TodaySummary, error) {
	// Build command to get today's summary using configurable Python executable
	cmd := exec.Command(s.config.PythonExecPath, "-m", "tt.cli", "today")
	cmd.Dir = s.config.TimeTrackerPath

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get today's summary: %s, output: %s", err.Error(), string(output))
	}

	// Parse the output to extract time information
	outputStr := string(output)

	// Debug: log the output
	fmt.Printf("DEBUG: Today summary output: %s\n", outputStr)

	// Simple parsing of the output
	totalHours := 0.0
	entryCount := 0
	var breakdown []Breakdown

	// Use regex to find patterns more reliably
	// Look for "Total Time: X.X hours" pattern
	totalTimeRegex := `Total Time:\s*([0-9.]+)\s*hours`
	if matches := regexp.MustCompile(totalTimeRegex).FindStringSubmatch(outputStr); len(matches) > 1 {
		if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
			totalHours = val
			fmt.Printf("DEBUG: Parsed total hours: %f\n", totalHours)
		}
	}

	// Look for breakdown entries with pattern "• Client/Project: X.Xh"
	breakdownRegex := `•\s*([^:]+):\s*([0-9.]+)h`
	matches := regexp.MustCompile(breakdownRegex).FindAllStringSubmatch(outputStr, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			clientProject := strings.TrimSpace(match[1])
			timeStr := strings.TrimSpace(match[2])
			if val, err := strconv.ParseFloat(timeStr, 64); err == nil {
				breakdown = append(breakdown, Breakdown{
					ClientProject: clientProject,
					Hours:         val,
					Minutes:       int(val * 60),
				})
				entryCount++
				fmt.Printf("DEBUG: Added breakdown entry: %s - %f hours\n", clientProject, val)
			}
		}
	}

	return &TodaySummary{
		TotalHours:   totalHours,
		TotalMinutes: int(totalHours * 60),
		EntryCount:   entryCount,
		Breakdown:    breakdown,
		RawOutput:    outputStr,
	}, nil
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
