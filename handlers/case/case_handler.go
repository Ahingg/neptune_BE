package caseHandler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	caseService "neptune/backend/services/case"
	"net/http"
	"time"
)

type CaseHandler struct {
	caseService caseService.CaseService
}

func NewCaseHandler(caseService caseService.CaseService) *CaseHandler {
	return &CaseHandler{caseService: caseService}
}

// CreateCase handles POST /api/cases
func (h *CaseHandler) CreateCase(c *gin.Context) {
	var req requests.CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.caseService.CreateCase(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create case: %v", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetCaseByID handles GET /api/cases/:caseId
func (h *CaseHandler) GetCaseByID(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("caseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.caseService.GetCaseByID(ctx, caseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve case: %v", err.Error())})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetAllCases handles GET /api/cases
func (h *CaseHandler) GetAllCases(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.caseService.GetAllCases(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve cases: %v", err.Error())})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateCase handles PUT /api/cases/:caseId
func (h *CaseHandler) UpdateCase(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("caseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID format"})
		return
	}
	var req requests.UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.caseService.UpdateCase(ctx, caseID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update case: %v", err.Error())})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteCase handles DELETE /api/cases/:caseId
func (h *CaseHandler) DeleteCase(c *gin.Context) {
	caseID, err := uuid.Parse(c.Param("caseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.caseService.DeleteCase(ctx, caseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete case: %v", err.Error())})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
