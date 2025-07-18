package contestHandler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"neptune/backend/pkg/requests"
	contestService "neptune/backend/services/contest"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContestHandler struct {
	contestService contestService.ContestService
}

func NewContestHandler(contestService contestService.ContestService) *ContestHandler {
	return &ContestHandler{contestService: contestService}
}

// CreateContest handles POST /api/contests
func (h *ContestHandler) CreateContest(c *gin.Context) {
	var req requests.CreateContestRequest

	// Log the raw request body for debugging
	body, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	fmt.Printf("Raw request body: %s\n", string(body))

	// Reset the request body so it can be read again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Parsed request: %+v\n", req)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.contestService.CreateContest(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create contest: %v", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetContestByID handles GET /api/contests/:contestId
func (h *ContestHandler) GetContestByID(c *gin.Context) {
	contestID, err := uuid.Parse(c.Param("contestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.contestService.GetContestByID(ctx, contestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve contest: %v", err.Error())})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetAllContests handles GET /api/contests
func (h *ContestHandler) GetAllContests(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.contestService.GetAllContests(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve contests: %v", err.Error())})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateContest handles PUT /api/contests/:contestId
func (h *ContestHandler) UpdateContest(c *gin.Context) {
	contestID, err := uuid.Parse(c.Param("contestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}
	var req requests.UpdateContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.contestService.UpdateContest(ctx, contestID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update contest: %v", err.Error())})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteContest handles DELETE /api/contests/:contestId
func (h *ContestHandler) DeleteContest(c *gin.Context) {
	contestID, err := uuid.Parse(c.Param("contestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.contestService.DeleteContest(ctx, contestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete contest: %v", err.Error())})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AddCasesToContest handles POST /api/contests/:contestId/cases
func (h *ContestHandler) AddCasesToContest(c *gin.Context) {
	contestID, err := uuid.Parse(c.Param("contestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}
	var req requests.AddCasesToContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.contestService.AddCasesToContest(ctx, contestID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add cases to contest: %v", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cases added to contest successfully"})
}

// AssignContestToClass handles POST /api/classes/:classTransactionId/contests
func (h *ContestHandler) AssignContestToClass(c *gin.Context) {
	classTransactionID, err := uuid.Parse(c.Param("classTransactionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class transaction ID format"})
		return
	}
	var req requests.AssignContestToClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.contestService.AssignContestToClass(ctx, classTransactionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to assign contest to class: %v", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetContestsForClass handles GET /api/classes/:classTransactionId/contests
func (h *ContestHandler) GetContestsForClass(c *gin.Context) {
	classTransactionID, err := uuid.Parse(c.Param("classTransactionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class transaction ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.contestService.GetContestsForClass(ctx, classTransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve contests for class: %v", err.Error())})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RemoveContestFromClass handles DELETE /api/classes/:classTransactionId/contests/:contestId
func (h *ContestHandler) RemoveContestFromClass(c *gin.Context) {
	classTransactionID, err := uuid.Parse(c.Param("classTransactionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class transaction ID format"})
		return
	}
	contestID, err := uuid.Parse(c.Param("contestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.contestService.RemoveContestFromClass(ctx, classTransactionID, contestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove contest from class: %v", err.Error())})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
