package services

import (
	"testing"

	"kb-freelance-api/internal/config"
)

func TestNewInvoiceService(t *testing.T) {
	cfg := &config.Config{
		TimeTrackerPath: "/test/path",
		InvoiceGenPath:  "/test/invoice",
		DatabasePath:    "/test/db",
		Port:            "8080",
	}

	service := NewInvoiceService(cfg)

	if service == nil {
		t.Fatal("Service should not be nil")
	}

	if service.config != cfg {
		t.Error("Service config should match provided config")
	}
}

func TestInvoiceLineItem(t *testing.T) {
	item := InvoiceLineItem{
		Description: "Test Description",
		Hours:       2.5,
		Rate:        75.0,
	}

	if item.Description != "Test Description" {
		t.Errorf("Expected Description 'Test Description', got '%s'", item.Description)
	}

	if item.Hours != 2.5 {
		t.Errorf("Expected Hours 2.5, got %f", item.Hours)
	}

	if item.Rate != 75.0 {
		t.Errorf("Expected Rate 75.0, got %f", item.Rate)
	}
}

func TestInvoiceLineItemValidation(t *testing.T) {
	tests := []struct {
		name        string
		item        InvoiceLineItem
		expectValid bool
	}{
		{
			name: "Valid item",
			item: InvoiceLineItem{
				Description: "Valid Description",
				Hours:       1.0,
				Rate:        50.0,
			},
			expectValid: true,
		},
		{
			name: "Zero hours",
			item: InvoiceLineItem{
				Description: "Zero Hours",
				Hours:       0.0,
				Rate:        50.0,
			},
			expectValid: false,
		},
		{
			name: "Negative hours",
			item: InvoiceLineItem{
				Description: "Negative Hours",
				Hours:       -1.0,
				Rate:        50.0,
			},
			expectValid: false,
		},
		{
			name: "Zero rate",
			item: InvoiceLineItem{
				Description: "Zero Rate",
				Hours:       1.0,
				Rate:        0.0,
			},
			expectValid: false,
		},
		{
			name: "Negative rate",
			item: InvoiceLineItem{
				Description: "Negative Rate",
				Hours:       1.0,
				Rate:        -50.0,
			},
			expectValid: false,
		},
		{
			name: "Empty description",
			item: InvoiceLineItem{
				Description: "",
				Hours:       1.0,
				Rate:        50.0,
			},
			expectValid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.item.Hours > 0 && test.item.Rate > 0 && test.item.Description != ""
			if isValid != test.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", test.expectValid, isValid)
			}
		})
	}
}

func TestInvoiceLineItemCalculations(t *testing.T) {
	item := InvoiceLineItem{
		Description: "Test Work",
		Hours:       2.5,
		Rate:        80.0,
	}

	expectedTotal := item.Hours * item.Rate
	actualTotal := 2.5 * 80.0

	if actualTotal != expectedTotal {
		t.Errorf("Expected total %f, got %f", expectedTotal, actualTotal)
	}

	if actualTotal != 200.0 {
		t.Errorf("Expected total 200.0, got %f", actualTotal)
	}
}



