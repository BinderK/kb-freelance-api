package services

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"kb-freelance-api/internal/config"
)

func TestNewTimeTrackerService(t *testing.T) {
	cfg := &config.Config{
		TimeTrackerPath: "/test/path",
		InvoiceGenPath:  "/test/invoice",
		DatabasePath:    "/test/db",
		Port:            "8080",
	}

	service := NewTimeTrackerService(cfg)

	if service == nil {
		t.Fatal("Service should not be nil")
	}

	if service.config != cfg {
		t.Error("Service config should match provided config")
	}
}

func TestTimerStatus(t *testing.T) {
	status := TimerStatus{
		IsRunning:       true,
		Client:          "Test Client",
		Project:         "Test Project",
		Description:     "Test Description",
		StartTime:       time.Now(),
		DurationMinutes: 30,
	}

	if !status.IsRunning {
		t.Error("IsRunning should be true")
	}

	if status.Client != "Test Client" {
		t.Errorf("Expected Client 'Test Client', got '%s'", status.Client)
	}

	if status.Project != "Test Project" {
		t.Errorf("Expected Project 'Test Project', got '%s'", status.Project)
	}
}

func TestTimeEntry(t *testing.T) {
	now := time.Now()
	endTime := now.Add(time.Hour)
	entry := TimeEntry{
		ID:              1,
		Client:          "Test Client",
		Project:         "Test Project",
		Description:     "Test Description",
		StartTime:       now,
		EndTime:         &endTime,
		DurationMinutes: 60,
	}

	if entry.ID != 1 {
		t.Errorf("Expected ID 1, got %d", entry.ID)
	}

	if entry.DurationMinutes != 60 {
		t.Errorf("Expected DurationMinutes 60, got %d", entry.DurationMinutes)
	}
}

func TestTodaySummary(t *testing.T) {
	breakdown := []Breakdown{
		{
			ClientProject: "Client A/Project A",
			Hours:         2.5,
			Minutes:       150,
		},
		{
			ClientProject: "Client B/Project B",
			Hours:         1.0,
			Minutes:       60,
		},
	}

	summary := TodaySummary{
		TotalHours:   3.5,
		TotalMinutes: 210,
		EntryCount:   2,
		Breakdown:    breakdown,
	}

	if summary.TotalHours != 3.5 {
		t.Errorf("Expected TotalHours 3.5, got %f", summary.TotalHours)
	}

	if summary.TotalMinutes != 210 {
		t.Errorf("Expected TotalMinutes 210, got %d", summary.TotalMinutes)
	}

	if summary.EntryCount != 2 {
		t.Errorf("Expected EntryCount 2, got %d", summary.EntryCount)
	}

	if len(summary.Breakdown) != 2 {
		t.Errorf("Expected 2 breakdown entries, got %d", len(summary.Breakdown))
	}
}

func TestBreakdown(t *testing.T) {
	breakdown := Breakdown{
		ClientProject: "Test Client/Test Project",
		Hours:         1.5,
		Minutes:       90,
	}

	if breakdown.ClientProject != "Test Client/Test Project" {
		t.Errorf("Expected ClientProject 'Test Client/Test Project', got '%s'", breakdown.ClientProject)
	}

	if breakdown.Hours != 1.5 {
		t.Errorf("Expected Hours 1.5, got %f", breakdown.Hours)
	}

	if breakdown.Minutes != 90 {
		t.Errorf("Expected Minutes 90, got %d", breakdown.Minutes)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "test", false},
		{"", "test", false},
		{"test", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		result := contains(test.s, test.substr)
		if result != test.expected {
			t.Errorf("contains(%q, %q) = %v, expected %v", test.s, test.substr, result, test.expected)
		}
	}
}

// Test the parsing logic with sample output
func TestParseTodaySummary(t *testing.T) {
	// Sample output from the Python CLI
	sampleOutput := `╭───────────────────────────── 📊 Today's Summary ─────────────────────────────╮
│ Total Time: 2.5 hours                                                        │
│                                                                              │
│ Breakdown:                                                                   │
│   • Client A/Project A: 1.5h                                                 │
│   • Client B/Project B: 1.0h                                                 │
│                                                                              │
╰──────────────────────────────────────────────────────────────────────────────╯`

	// Test the regex patterns
	totalTimeRegex := `Total Time:\s*([0-9.]+)\s*hours`
	matches := regexp.MustCompile(totalTimeRegex).FindStringSubmatch(sampleOutput)
	if len(matches) < 2 {
		t.Fatal("Should find total time match")
	}

	totalHours, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		t.Fatalf("Failed to parse total hours: %v", err)
	}

	if totalHours != 2.5 {
		t.Errorf("Expected total hours 2.5, got %f", totalHours)
	}

	// Test breakdown parsing
	breakdownRegex := `•\s*([^:]+):\s*([0-9.]+)h`
	breakdownMatches := regexp.MustCompile(breakdownRegex).FindAllStringSubmatch(sampleOutput, -1)
	if len(breakdownMatches) != 2 {
		t.Errorf("Expected 2 breakdown matches, got %d", len(breakdownMatches))
	}

	// Check first breakdown entry
	if len(breakdownMatches[0]) >= 3 {
		clientProject := strings.TrimSpace(breakdownMatches[0][1])
		timeStr := strings.TrimSpace(breakdownMatches[0][2])

		if clientProject != "Client A/Project A" {
			t.Errorf("Expected client project 'Client A/Project A', got '%s'", clientProject)
		}

		hours, err := strconv.ParseFloat(timeStr, 64)
		if err != nil {
			t.Fatalf("Failed to parse hours: %v", err)
		}

		if hours != 1.5 {
			t.Errorf("Expected 1.5 hours, got %f", hours)
		}
	}
}
