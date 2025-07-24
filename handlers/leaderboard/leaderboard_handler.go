package leaderboardHand

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	contestService "neptune/backend/services/contest"
	"neptune/backend/services/leaderboard"
	"net/http"
)

type LeaderboardHandler struct {
	service     leaderboardServ.Service
	contestServ contestService.ContestService
}

func NewLeaderboardHandler(service leaderboardServ.Service, contestServ contestService.ContestService) *LeaderboardHandler {
	return &LeaderboardHandler{
		service:     service,
		contestServ: contestServ,
	}
}

func (h *LeaderboardHandler) GetGlobalContestLeaderboard(c *gin.Context) {
	contestIDStr := c.Param("contestId")

	contestID, err := uuid.Parse(contestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}

	leaderboardData, err := h.service.GetGlobalContestLeaderboard(c.Request.Context(), contestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate leaderboard", "details": err.Error()})
		return
	}

	contestCases, err := h.contestServ.GetContestCases(c.Request.Context(), contestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contest cases", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contest_id":  contestIDStr,
		"cases":       contestCases,
		"leaderboard": leaderboardData,
	})
}

func (h *LeaderboardHandler) GetClassContestLeaderboard(c *gin.Context) {
	classIDStr := c.Param("classTransactionId")
	contestIDStr := c.Param("contestId")

	fmt.Println(classIDStr)

	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID format"})
		return
	}

	contestID, err := uuid.Parse(contestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID format"})
		return
	}

	leaderboardData, err := h.service.GetContestLeaderboard(c.Request.Context(), classID, contestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate leaderboard", "details": err.Error()})
		return
	}

	// get case
	contestCases, err := h.contestServ.GetContestCases(c.Request.Context(), contestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contest cases", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class_transaction_id": classIDStr,
		"contest_id":           contestIDStr,
		"cases":                contestCases,
		"leaderboard":          leaderboardData,
	})
}
