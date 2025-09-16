package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
	// For now, we'll use the first line item only since the Python CLI doesn't support multiple items
	// TODO: Implement proper multiple line items support
	if len(lineItems) == 0 {
		return nil, fmt.Errorf("at least one line item is required")
	}

	// Use the first line item
	item := lineItems[0]

	// Build command to generate invoice using the original Python CLI
	cmd := exec.Command(s.config.PythonExecPath, "-m", "src.main",
		"-c", clientName,
		"-e", clientEmail,
		"-d", item.Description,
		"-h", fmt.Sprintf("%.2f", item.Hours),
		"-r", fmt.Sprintf("%.2f", item.Rate),
	)

	// Set working directory to invoice generator path
	cmd.Dir = s.config.InvoiceGenPath

	// Add optional parameters
	// Always provide notes parameter to avoid interactive prompts
	notesValue := notes
	if notesValue == "" {
		notesValue = "Generated via API"
	}
	cmd.Args = append(cmd.Args, "--notes", notesValue)
	if date != "" {
		cmd.Args = append(cmd.Args, "--date", date)
	}

	// Execute command
	fmt.Printf("DEBUG: Running invoice command: %v\n", cmd.Args)
	fmt.Printf("DEBUG: Working directory: %s\n", cmd.Dir)
	fmt.Printf("DEBUG: Python executable: %s\n", s.config.PythonExecPath)

	// Check if the invoice generator directory exists
	if _, err := os.Stat(s.config.InvoiceGenPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("invoice generator path does not exist: %s", s.config.InvoiceGenPath)
	}

	// Clean up any existing PDF files to avoid conflicts BEFORE running the Python script
	outputDir := filepath.Join(s.config.InvoiceGenPath, "output")
	if files, err := os.ReadDir(outputDir); err == nil {
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".pdf" {
				oldPdfPath := filepath.Join(outputDir, file.Name())
				fmt.Printf("DEBUG: Removing old PDF file: %s\n", oldPdfPath)
				os.Remove(oldPdfPath)
			}
		}
	}

	// First, test if Python is working
	testCmd := exec.Command(s.config.PythonExecPath, "--version")
	testOutput, testErr := testCmd.CombinedOutput()
	if testErr != nil {
		return nil, fmt.Errorf("Python executable not working: %s, output: %s", testErr.Error(), string(testOutput))
	}
	fmt.Printf("DEBUG: Python test successful: %s", string(testOutput))

	output, err := cmd.CombinedOutput()
	fmt.Printf("DEBUG: Command output: %s\n", string(output))
	if err != nil {
		fmt.Printf("DEBUG: Command error: %v\n", err)
		fmt.Printf("DEBUG: Command exit code: %d\n", cmd.ProcessState.ExitCode())

		// Check if the error is due to interactive prompts
		if strings.Contains(string(output), "Aborted!") {
			return nil, fmt.Errorf("invoice generation failed due to interactive prompts. This usually means the Python script is expecting user input. Output: %s", string(output))
		}

		return nil, fmt.Errorf("failed to generate invoice: %s, output: %s", err.Error(), string(output))
	}

	// Generate a unique filename for this invoice
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("invoice_%s_%s.pdf", clientName, timestamp)
	// Replace spaces and special characters in filename
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "_")

	// The PDF should be generated in the output directory
	pdfPath := filepath.Join(outputDir, "invoice.pdf")
	fmt.Printf("DEBUG: Looking for PDF at: %s\n", pdfPath)

	// Check if the PDF was actually created
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		// List files in output directory for debugging
		if files, err := os.ReadDir(outputDir); err == nil {
			fmt.Printf("DEBUG: Files in output directory: %v\n", files)
		}

		// Try to find any PDF file in the output directory
		if files, err := os.ReadDir(outputDir); err == nil {
			for _, file := range files {
				if filepath.Ext(file.Name()) == ".pdf" {
					fmt.Printf("DEBUG: Found PDF file: %s\n", file.Name())
					pdfPath = filepath.Join(outputDir, file.Name())
					break
				}
			}
		}

		// If still no PDF found, return error with more details
		if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("PDF file was not created at %s. Python command may have failed. Check server logs for details.", pdfPath)
		}
	}

	return map[string]interface{}{
		"status":       "success",
		"message":      "Invoice generated successfully",
		"pdf_path":     pdfPath,
		"filename":     filename,
		"download_url": "/files/invoice.pdf",
		"output":       string(output),
	}, nil
}
