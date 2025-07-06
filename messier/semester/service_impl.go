package semester

import (
	"context"
	"fmt"
	"neptune/backend/messier"
	"os"
)

type semesterServiceImpl struct {
}

func NewExternalSemesterService() MessierSemesterService {
	return &semesterServiceImpl{}
}

func (s semesterServiceImpl) GetSemesters(ctx context.Context, authToken string) ([]GetSemestersResponse, error) {
	var result []GetSemestersResponse

	url := fmt.Sprintf("%s/Semester/GetSemestersWithActiveDate", os.Getenv("MESSIER_API_URL"))
	err := messier.SendRequest(ctx, "GET", url, nil, &result, authToken)
	if err != nil {
		return nil, err
	}

	return result, nil
}
