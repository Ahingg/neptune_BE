package amqp_messages

import (
	"github.com/google/uuid"
	submissionModel "neptune/backend/models/submission"
)

type JudgeQueueMessage struct {
	SubmissionID uuid.UUID `json:"submission_id" binding:"required"`
}

type ResultQueueMessage struct {
	SubmissionID uuid.UUID                          `json:"submission_id" binding:"required"`
	FinalStatus  submissionModel.SubmissionStatus   `json:"final_status" binding:"required"`
	Score        int                                `json:"score" binding:"required"`
	Results      []submissionModel.SubmissionResult `json:"results" binding:"required"`
}
