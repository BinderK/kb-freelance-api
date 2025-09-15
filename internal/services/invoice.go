package services

import (
	"fmt"
	"os/exec"

	"kb-freelance-api/internal/config"
)

type InvoiceService struct {
	config *config.Config
}

func NewInvoiceService(cfg *config.Config) *InvoiceService {
	return &InvoiceService{config: cfg}
}

type InvoiceLineItem struct {
	Description string  `json:"description"`
	Hours       float64 `json:"hours"`
	Rate        float64 `json:"rate"`
}

type InvoiceRequest struct {
	ClientName  string            `json:"client_name"`
	ClientEmail string            `json:"client_email"`
	LineItems   []InvoiceLineItem `json:"line_items"`
	Notes       string            `json:"notes"`
	Date        string            `json:"date"`
}

func (s *InvoiceService) GenerateInvoice(clientName, clientEmail string, lineItems []InvoiceLineItem, notes, date string) (map[string]interface{}, error) {
	// Build command to generate invoice using configurable Python executable
	cmd := exec.Command(s.config.PythonExecPath, "-m", "src.main",
		"-c", clientName,
		"-e", clientEmail,
	)

	// Set working directory to invoice generator path
	cmd.Dir = s.config.InvoiceGenPath

	// Add line items as arguments
	// For now, we'll use the first line item as a simple example
	// In a real implementation, you'd need to handle multiple line items
	if len(lineItems) > 0 {
		item := lineItems[0]
		cmd.Args = append(cmd.Args,
			"-d", item.Description,
			"-h", fmt.Sprintf("%.2f", item.Hours),
			"-r", fmt.Sprintf("%.2f", item.Rate),
		)
	}

	// Add optional parameters
	if notes != "" {
		cmd.Args = append(cmd.Args, "--notes", notes)
	}
	if date != "" {
		cmd.Args = append(cmd.Args, "--date", date)
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice: %s, output: %s", err.Error(), string(output))
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Invoice generated successfully",
		"output":  string(output),
	}, nil
}
