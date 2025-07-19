package testCaseHand

import (
	"context"
	"github.com/gin-gonic/gin"
	"neptune/backend/pkg/requests"
	caseService "neptune/backend/services/case"
	testCaseServ "neptune/backend/services/test_case"
	"time"
)

type TestCaseHandler struct {
	testCaseService testCaseServ.TestCaseService
	caseServ        caseService.CaseService
}

func NewTestCaseHandler(testCaseService testCaseServ.TestCaseService, caseServ caseService.CaseService) *TestCaseHandler {
	return &TestCaseHandler{
		testCaseService: testCaseService,
		caseServ:        caseServ,
	}
}

func (h *TestCaseHandler) UploadTestCasesHandler(c *gin.Context) {
	var req requests.AddTestCaseRequest
	if err := req.ParseFormData(c); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	err := h.testCaseService.UploadTestCases(ctx, req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload test cases", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Test cases uploaded successfully"})
}

func (h *TestCaseHandler) GetTestCasesByCaseIDHandler(c *gin.Context) {
	caseID := c.Param("case_id")
	if caseID == "" {
		c.JSON(400, gin.H{"error": "Case ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	testCases, err := h.testCaseService.GetTestCasesByCaseID(ctx, caseID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve test cases", "details": err.Error()})
		return
	}

	c.JSON(200, testCases)
}
