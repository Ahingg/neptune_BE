package requests

type CreateCaseRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PDFFileUrl    string `json:"pdf_file_url"`
	TimeLimitMs   int    `json:"time_limit_ms" binding:"required,min=1"`
	MemoryLimitMb int    `json:"memory_limit_mb" binding:"required,min=1"`
}

type UpdateCaseRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PDFFileUrl    string `json:"pdf_file_url"`
	TimeLimitMs   int    `json:"time_limit_ms" binding:"required,min=1"`
	MemoryLimitMb int    `json:"memory_limit_mb" binding:"required,min=1"`
}
