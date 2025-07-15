package semester

import (
	"fmt"
	"neptune/backend/pkg/responses"
	"neptune/backend/services/internal_semester"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SemesterHandler struct {
	internalSemesterService internal_semester.SemesterService
}

func NewSemesterHandler(internalSemesterService internal_semester.SemesterService) *SemesterHandler {
	return &SemesterHandler{
		internalSemesterService: internalSemesterService,
	}
}

func (h *SemesterHandler) SyncSemestersHandler(c *gin.Context) {
	// Get the UserID from the Gin context (set by AuthRequired middleware)
	requestMakerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	requestMakerIDStr, ok := requestMakerID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type in context"})
		return
	}

	err := h.internalSemesterService.SyncSemester(c.Request.Context(), requestMakerIDStr)
	if err != nil {
		// Differentiate between auth/permission errors and internal errors
		if err.Error() == fmt.Sprintf("unauthorized: no Messier token found for user %s", requestMakerIDStr) ||
			err.Error() == fmt.Sprintf("unauthorized: Messier token for user %s has expired", requestMakerIDStr) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sync semesters: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semesters synchronized successfully"})
}

// GetSemestersHandler retrieves all semesters from the internal database.
// Accessible by authenticated users with appropriate roles (e.g., student, lecturer, admin).
func (h *SemesterHandler) GetSemestersHandler(c *gin.Context) {
	semesters, err := h.internalSemesterService.GetInternalSemesters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve semesters: %v", err)})
		return
	}

	fmt.Printf("Found %d semesters in database\n", len(semesters))
	for i, s := range semesters {
		fmt.Printf("Semester %d: ID=%s, Description=%s\n", i+1, s.ID, s.Description)
	}

	// Transform model.Semester to your desired GetSemestersResponse format
	// if needed, otherwise just send model.Semester array
	responseSemesters := make([]responses.SemesterResponse, len(semesters))

	for i, s := range semesters {
		responseSemesters[i] = responses.SemesterResponse{
			Description: s.Description,
			End:         s.End, // Direct assignment of *time.Time (will be null if DB is null)
			SemesterID:  s.ID,
			Start:       s.Start,
		}
	}

	fmt.Printf("Sending %d semesters to frontend\n", len(responseSemesters))
	c.JSON(http.StatusOK, responseSemesters)
}

// DebugSemestersHandler is a temporary debug endpoint to check database state
func (h *SemesterHandler) DebugSemestersHandler(c *gin.Context) {
	// Get all semesters
	semesters, err := h.internalSemesterService.GetInternalSemesters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to retrieve semesters: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":     len(semesters),
		"semesters": semesters,
		"message":   "Debug info for semesters",
	})
}
