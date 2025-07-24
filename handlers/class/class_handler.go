package class

import (
	"context"
	"neptune/backend/pkg/requests"
	"neptune/backend/services/internal_class"
	"time"

	"neptune/backend/messier/constants"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	internalClassService internal_class.ClassService
}

func NewClassHandler(internalClassService internal_class.ClassService) *ClassHandler {
	return &ClassHandler{
		internalClassService: internalClassService,
	}
}

func (h *ClassHandler) SyncClassesHandler(c *gin.Context) {
	requestMakerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "user_id not found in context"})
		return
	}

	semesterId := c.Query("semester_id")
	if semesterId == "" {
		c.JSON(400, gin.H{"error": "invalid param"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	err := h.internalClassService.SyncClasses(ctx, semesterId, requestMakerID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to sync classes", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "classes synced successfully"})
}

func (h *ClassHandler) SyncClassStudentsHandler(c *gin.Context) {
	requestMakerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "user_id not found in context"})
		return
	}

	var req requests.SyncClassStudentAndAssistantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	err := h.internalClassService.SyncClassStudents(ctx, req.SemesterID, req.CourseID, requestMakerID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to sync class students", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "class students synced successfully"})
}

func (h *ClassHandler) SyncClassAssistantsHandler(c *gin.Context) {
	requestMakerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "user_id not found in context"})
		return
	}

	var req requests.SyncClassStudentAndAssistantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	err := h.internalClassService.SyncClassAssistants(ctx, req.SemesterID, req.CourseID, requestMakerID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to sync class assistants", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "class assistants synced successfully"})
}

func (h *ClassHandler) GetClassesBySemesterAndCourseHandler(c *gin.Context) {
	semesterID := c.Query("semester_id")
	if semesterID == "" {
		c.JSON(400, gin.H{"error": "invalid semester_id"})
		return
	}

	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(400, gin.H{"error": "invalid course_id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	classes, err := h.internalClassService.GetClassesBySemesterAndCourse(ctx, semesterID, courseID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get classes", "details": err.Error()})
		return
	}
	c.JSON(200, classes)
}

func (h *ClassHandler) GetClassDetailBySemesterAndCourseHandler(c *gin.Context) {
	semesterID := c.Query("semester_id")
	if semesterID == "" {
		c.JSON(400, gin.H{"error": "invalid semester_id"})
		return
	}

	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(400, gin.H{"error": "invalid course_id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	classDetails, err := h.internalClassService.GetClassDetailBySemesterAndCourse(ctx, semesterID, courseID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get class details", "details": err.Error()})
		return
	}
	c.JSON(200, classDetails)
}

func (h *ClassHandler) GetClassDetailBySemesterCourseAndStudentHandler(c *gin.Context) {
	semesterID := c.Query("semester_id")
	if semesterID == "" {
		c.JSON(400, gin.H{"error": "invalid semester_id"})
		return
	}

	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(400, gin.H{"error": "invalid course_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "user_id not found in context"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	classDetails, err := h.internalClassService.GetClassDetailBySemesterCourseAndStudent(ctx, semesterID, courseID, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get class details", "details": err.Error()})
		return
	}
	c.JSON(200, classDetails)
}

func (h *ClassHandler) GetClassDetailByTransactionIDHandler(c *gin.Context) {
	classTransactionID := c.Query("class_transaction_id")
	if classTransactionID == "" {
		c.JSON(400, gin.H{"error": "invalid class_transaction_id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	classDetail, err := h.internalClassService.GetClassDetailByTransactionID(ctx, classTransactionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get class detail", "details": err.Error()})
		return
	}
	if classDetail == nil {
		c.JSON(404, gin.H{"error": "class not found"})
		return
	}
	c.JSON(200, classDetail)
}

type Course struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetCoursesHandler(c *gin.Context) {
	courses := []Course{
		{ID: constants.AlgoprogID1, Name: "COMP6047001"},
		{ID: constants.AlgoprogID2, Name: "COMP6878051"},
	}
	c.JSON(200, courses)
}
