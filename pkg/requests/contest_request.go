package requests

type CreateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
