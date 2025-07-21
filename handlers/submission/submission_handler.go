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
	// 1. Create request and let it parse itself from the context
	var req requests.SubmitCodeRequest
	if err := req.ParseAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Get userID from the auth middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uId, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 3. Call the service with the parsed and validated request
	submission, err := h.service.SubmitCode(c.Request.Context(), &req, uId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, submission)
}
