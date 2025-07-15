package caseHandler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"neptune/backend/pkg/requests"
	caseService "neptune/backend/services/case"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse form data: %v", err.Error())})
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	timeLimitMsStr := c.PostForm("time_limit_ms")
	memoryLimitMbStr := c.PostForm("memory_limit_mb")

	timeLimitMs, err := strconv.Atoi(timeLimitMsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_limit_ms format"})
		return
	}

	memoryLimitMb, err := strconv.Atoi(memoryLimitMbStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid memory_limit_mb format"})
		return
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if timeLimitMs <= 0 || memoryLimitMb <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time and memory limits must be positive"})
		return
	}

	file, err := c.FormFile("pdf_file") // "pdf_file" is the name of the input field in the form
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to get PDF file: %v", err.Error())})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	uniqueFilename := uuid.New().String() + ext

	PDF_UPLOAD := "./public/case_file"
	filePath := filepath.Join(PDF_UPLOAD, uniqueFilename)
	fileURL := "/public/case_file/" + uniqueFilename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving uploaded file %s: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PDF file"})
		return
	}

	serviceReq := requests.CreateCaseRequest{
		Name:          name,
		Description:   description,
		TimeLimitMs:   timeLimitMs,
		MemoryLimitMb: memoryLimitMb,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	resp, err := h.caseService.CreateCase(ctx, serviceReq, fileURL)
	if err != nil {
		os.Remove(filePath)
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
