package submissionHand

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	submissionServ "neptune/backend/services/submission"
	"net/http"
)

type SubmissionHandler struct {
	service submissionServ.SubmissionService
}

func NewSubmissionHandler(service submissionServ.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

func (h *SubmissionHandler) SubmitCode(c *gin.Context) {
	var req requests.SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from the auth middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	submission, err := h.service.SubmitCode(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Status 202 Accepted is appropriate here, as the request has been
	// accepted for processing, but the processing is not yet complete.
	c.JSON(http.StatusAccepted, submission)
}
