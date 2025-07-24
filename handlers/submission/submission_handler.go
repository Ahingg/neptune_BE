package submissionHand

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
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

func (h *SubmissionHandler) GetSubmissionByUserInContest(c *gin.Context) {
	userIdStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user_id in token"})
		return
	}

	userId, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId does not satisfy the requirement"})
		return
	}

	contestIdStr := c.Param("contestId")
	contestId, err := uuid.Parse(contestIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid type for contestId"})
		return
	}

	classIdStr := c.Query("class_transaction_id")
	var resp []responses.GetUserSubmissionsResponse

	if classIdStr == "" {
		resp, err = h.service.GetSubmissionByUserInContest(c.Request.Context(), userId, contestId, nil)
	} else {
		classId, err := uuid.Parse(classIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid type for classId"})
			return
		}
		fmt.Println(classId, contestId, userId)
		resp, err = h.service.GetSubmissionByUserInContest(c.Request.Context(), userId, contestId, &classId)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SubmissionHandler) GetClassContestSubmissions(c *gin.Context) {

	classIdStr := c.Query("class_transaction_id")
	if classIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_transaction_id is required"})
		return
	}

	classId, err := uuid.Parse(classIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid type for class_transaction_id"})
		return
	}

	contestIdStr := c.Param("contestId")
	contestId, err := uuid.Parse(contestIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid type for contestId"})
		return
	}

	resp, err := h.service.GetClassContestSubmissions(c.Request.Context(), classId, contestId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
