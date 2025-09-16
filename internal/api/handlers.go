package api

import (
	"net/http"
	"strconv"

	"kb-freelance-api/internal/services"

	"github.com/gin-gonic/gin"
)

// Time tracking handlers

type StartTimerRequest struct {
	Client      string `json:"client" binding:"required"`
	Project     string `json:"project" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) startTimer(c *gin.Context) {
	var req StartTimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result, err := s.timeTrackerService.StartTimer(req.Client, req.Project, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (s *Server) stopTimer(c *gin.Context) {
	result, err := s.timeTrackerService.StopTimer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (s *Server) getTimerStatus(c *gin.Context) {
	status, err := s.timeTrackerService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// If no timer is running, return null
	if status == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": status})
}

func (s *Server) getTimeEntries(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	entries, err := s.timeTrackerService.GetRecentEntries(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": entries})
}

func (s *Server) getTodaySummary(c *gin.Context) {
	summary, err := s.timeTrackerService.GetTodaySummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

// Invoice handlers

type GenerateInvoiceRequest struct {
	ClientName  string                   `json:"client_name" binding:"required"`
	ClientEmail string                   `json:"client_email" binding:"required"`
	LineItems   []InvoiceLineItemRequest `json:"line_items" binding:"required"`
	Notes       string                   `json:"notes"`
	Date        string                   `json:"date"`
}

type InvoiceLineItemRequest struct {
	Description string  `json:"description" binding:"required"`
	Hours       float64 `json:"hours" binding:"required"`
	Rate        float64 `json:"rate" binding:"required"`
}

func (s *Server) generateInvoice(c *gin.Context) {
	var req GenerateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Convert request line items to service line items
	lineItems := make([]services.InvoiceLineItem, len(req.LineItems))
	for i, item := range req.LineItems {
		lineItems[i] = services.InvoiceLineItem{
			Description: item.Description,
			Hours:       item.Hours,
			Rate:        item.Rate,
		}
	}

	result, err := s.invoiceService.GenerateInvoice(req.ClientName, req.ClientEmail, lineItems, req.Notes, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (s *Server) previewInvoice(c *gin.Context) {
	// TODO: Implement invoice preview
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Invoice preview not yet implemented"})
}
