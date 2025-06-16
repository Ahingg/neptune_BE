package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"neptune/backend/models"

	"github.com/gin-gonic/gin"
)

type Judge0Submission struct {
	SourceCode     string `json:"source_code"`
	LanguageID     int    `json:"language_id"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

type Judge0Response struct {
	Status struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
	Stderr string `json:"stderr"`
	Output string `json:"output"`
}

func CreateSubmission(c *gin.Context) {
	var req models.SubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the case and its test cases
	case_, err := models.GetCase(db, req.CaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	// Create submission record
	submission := models.Submission{
		UserID:   c.GetUint("user_id"),
		CaseID:   req.CaseID,
		Code:     req.Code,
		Language: req.Language,
		Status:   "Pending",
	}

	if err := models.CreateSubmission(db, &submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process each test case
	for _, testCase := range case_.TestCases {
		// Skip hidden test cases for students
		if testCase.IsHidden && c.GetString("role") == "Student" {
			continue
		}

		// Prepare Judge0 submission
		judge0Sub := Judge0Submission{
			SourceCode:     req.Code,
			LanguageID:     getLanguageID(req.Language),
			Stdin:          testCase.Input,
			ExpectedOutput: testCase.Output,
		}

		// Send to Judge0
		judge0Resp, err := sendToJudge0(judge0Sub)
		if err != nil {
			models.UpdateSubmissionStatus(db, submission.ID, "Error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process submission"})
			return
		}

		// Check result
		if judge0Resp.Status.ID != 3 { // 3 is Accepted
			models.UpdateSubmissionStatus(db, submission.ID, "Wrong Answer", judge0Resp.Stderr)
			c.JSON(http.StatusOK, gin.H{
				"status": "Wrong Answer",
				"error":  judge0Resp.Stderr,
			})
			return
		}
	}

	// All test cases passed
	models.UpdateSubmissionStatus(db, submission.ID, "Accepted", "")
	c.JSON(http.StatusOK, gin.H{"status": "Accepted"})
}

func GetSubmission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	submission, err := models.GetSubmission(db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	// Check if user has permission to view this submission
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if userID != submission.UserID && role != "Lecturer" && role != "Assistant" && role != "SubDev" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

func GetSubmissions(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var submissions []models.Submission
	var err error

	if role == "Student" {
		submissions, err = models.GetUserSubmissions(db, userID)
	} else {
		caseID, _ := strconv.ParseUint(c.Query("case_id"), 10, 32)
		if caseID > 0 {
			submissions, err = models.GetCaseSubmissions(db, uint(caseID))
		} else {
			// Get all submissions for staff
			submissions, err = models.GetAllSubmissions(db)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

func GetClassLeaderboard(c *gin.Context) {
	classID := c.Param("classId")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class ID is required"})
		return
	}

	leaderboard, err := models.GetClassLeaderboard(db, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

func getLanguageID(language string) int {
	switch language {
	case "C":
		return 50 // GCC 9.2.0
	case "Python":
		return 71 // Python 3.8.1
	default:
		return 0
	}
}

func sendToJudge0(submission Judge0Submission) (*Judge0Response, error) {
	jsonData, err := json.Marshal(submission)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/submissions", os.Getenv("JUDGE0_API_URL")),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var judge0Resp Judge0Response
	if err := json.NewDecoder(resp.Body).Decode(&judge0Resp); err != nil {
		return nil, err
	}

	return &judge0Resp, nil
}
