package semester

import (
	"context"
	"fmt"
	"neptune/backend/messier"
	"os"
	"strings"
)

type semesterServiceImpl struct {
}

func NewExternalSemesterService() MessierSemesterService {
	return &semesterServiceImpl{}
}

func (s semesterServiceImpl) GetSemesters(ctx context.Context, authToken string) ([]GetSemestersResponse, error) {
	var result []GetSemestersResponse

	baseURL := strings.TrimRight(os.Getenv("MESSIER_API_URL"), "/")
	url := fmt.Sprintf("%s/Semester/GetSemestersWithActiveDate", baseURL)
	fmt.Printf("Fetching semesters from: %s\n", url)
	err := messier.SendRequest(ctx, "GET", url, nil, &result, authToken)
	if err != nil {
		fmt.Printf("Error fetching semesters from Messier: %v\n", err)
		return nil, err
	}

	fmt.Printf("Received %d semesters from Messier API\n", len(result))
	for i, semester := range result {
		fmt.Printf("Messier Semester %d: ID=%s, Description=%s\n", i+1, semester.SemesterID, semester.Description)
	}

	return result, nil
}
