package handlers

import (
	"net/http"
	"strconv"

	"neptune/backend/models"

	"github.com/gin-gonic/gin"
)

func GetCases(c *gin.Context) {
	classID := c.Query("class_id")
	isPublic := c.Query("is_public") == "true"

	cases, err := models.GetCases(db, classID, isPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cases)
}

func GetCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	case_, err := models.GetCase(db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	c.JSON(http.StatusOK, case_)
}

func CreateCase(c *gin.Context) {
	var case_ models.Case
	if err := c.ShouldBindJSON(&case_); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has permission to create cases
	role := c.GetString("role")
	if role != "Assistant" && role != "SubDev" && role != "Lecturer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.CreateCase(db, &case_); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, case_)
}

func UpdateCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	var case_ models.Case
	if err := c.ShouldBindJSON(&case_); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	case_.ID = uint(id)

	// Check if user has permission to update cases
	role := c.GetString("role")
	if role != "Assistant" && role != "SubDev" && role != "Lecturer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.UpdateCase(db, &case_); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, case_)
}

func DeleteCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	// Check if user has permission to delete cases
	role := c.GetString("role")
	if role != "Assistant" && role != "SubDev" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DeleteCase(db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Case deleted successfully"})
}

func AddTestCase(c *gin.Context) {
	var testCase models.TestCase
	if err := c.ShouldBindJSON(&testCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has permission to add test cases
	role := c.GetString("role")
	if role != "Assistant" && role != "SubDev" && role != "Lecturer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.AddTestCase(db, &testCase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, testCase)
}

func DeleteTestCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test case ID"})
		return
	}

	// Check if user has permission to delete test cases
	role := c.GetString("role")
	if role != "Assistant" && role != "SubDev" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DeleteTestCase(db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test case deleted successfully"})
}
