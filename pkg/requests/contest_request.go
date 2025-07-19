package requests

type CreateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Scope       string `json:"scope" binding:"required"` // e.g., "public", "class"
}

type UpdateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Scope       string `json:"scope" binding:"required"` // e.g., "public", "class"
}
