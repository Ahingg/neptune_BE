package submissionHand

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	submissionServ "neptune/backend/services/submission"
	"net/http"
)

type SubmissionReviewHandler struct {
	service submissionServ.SubmissionReviewService
}

func NewSubmissionReviewHandler(service submissionServ.SubmissionReviewService) *SubmissionReviewHandler {
	return &SubmissionReviewHandler{service: service}
}

// ViewCode handles the request to get the submission code as plain text.
func (h *SubmissionReviewHandler) ViewCode(c *gin.Context) {
	submissionIDStr := c.Param("submissionId")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
		return
	}

	code, contentType, err := h.service.GetSubmissionCode(c.Request.Context(), submissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Use c.Data to send raw bytes with a specific content type.
	c.Data(http.StatusOK, contentType, code)
}

// DownloadCode handles the request to download the submission code as a zip file.
func (h *SubmissionReviewHandler) DownloadCode(c *gin.Context) {
	submissionIDStr := c.Param("submissionId")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
		return
	}

	zipData, downloadFilename, err := h.service.GetSubmissionCodeAsZip(c.Request.Context(), submissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Set headers to prompt the browser to download the file.
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename))
	c.Data(http.StatusOK, "application/zip", zipData.Bytes())
}
