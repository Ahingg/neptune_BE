package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
)

type AddTestCaseRequest struct {
	CaseID uuid.UUID
	File   *multipart.FileHeader
}

func (r *AddTestCaseRequest) ParseFormData(c *gin.Context) error {
	caseID, err := uuid.Parse(c.Param("case_id"))
	if err != nil {
		return err
	}
	r.CaseID = caseID

	file, err := c.FormFile("test_case_zip")
	if err != nil {
		return err
	}
	r.File = file

	return nil
}
