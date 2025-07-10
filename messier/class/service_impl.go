package class

import (
	"context"
	"fmt"
	"log"
	"neptune/backend/messier"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type messierClassService struct {
}

func (s *messierClassService) GetAssistantInitialFromStudentTransaction(ctx context.Context, studentNIM string,
	semesterId string, authToken string) ([]GetStudentClassTransactionWithAssistantResponse, error) {
	var response []GetStudentClassTransactionWithAssistantResponse
	endpointUrl := fmt.Sprintf("%s/Student/GetStudentClassTransactionWithAssistant?nim=%s&semesterId=%s",
		os.Getenv("MESSIER_API_URL"), studentNIM, semesterId)
	err := messier.SendRequest(ctx, "GET", endpointUrl, nil, &response,
		authToken)
	//fmt.Printf("StudentNim: %s, response: %+v\n", studentNIM, response)
	return response, err
}

func (s *messierClassService) GetAssistantDetailFromAssistantInitial(ctx context.Context,
	initial string, generation string, authToken string) (*GetAssistantDetailResponse, error) {

	var response []GetAssistantDetailResponse

	// Clean up base URL to avoid double slashes
	baseUrl := strings.TrimRight(os.Getenv("MESSIER_API_URL"), "/")

	// Escape query parameters safely
	escapedInitial := url.QueryEscape(initial)
	escapedGeneration := url.QueryEscape(generation)

	// Build the full endpoint URL
	endpointUrl := fmt.Sprintf("%s/Assistant?initial=%s&generation=%s",
		baseUrl, escapedInitial, escapedGeneration)

	log.Println("Calling Assistant detail API:", endpointUrl)

	// Perform the request
	err := messier.SendRequest(ctx, http.MethodGet, endpointUrl, nil, &response, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get assistant detail: %w", err)
	}

	return &response[0], nil
}

func (s *messierClassService) GetStudentFromClassTransaction(ctx context.Context,
	semesterId, courseOutlineId, className, authToken string) ([]GetClassStudentsResponse, error) {

	var students []GetClassStudentsResponse

	baseUrl := strings.TrimRight(os.Getenv("MESSIER_API_URL"), "/")
	escapedClassName := url.QueryEscape(className)

	endpointUrl := fmt.Sprintf("%s/Student/Class?coId=%s&className=%s&semesterId=%s",
		baseUrl, courseOutlineId, escapedClassName, semesterId)

	log.Println("Calling URL:", endpointUrl)

	err := messier.SendRequest(ctx, http.MethodGet, endpointUrl, nil, &students, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get class students: %w", err)
	}

	return students, nil
}

func (s *messierClassService) GetClassesBySemesterAndCourseOutline(ctx context.Context, semesterId string,
	courseId string, authToken string) ([]GetClassBySemesterAndCourseResponse, error) {
	var result []GetClassBySemesterAndCourseResponse
	// bakal dapetin semua kelas yang ada di semester dan course itu.
	endpointUrl := fmt.Sprintf("%s/ClassTransaction/GetClassBySemesterAndCourseOutline?semesterId=%s&courseOutlineId=%s",
		os.Getenv("MESSIER_API_URL"),
		semesterId, courseId)
	err := messier.SendRequest(ctx, "GET", endpointUrl, nil, &result, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}
	return result, nil
}

func NewMessierClassService() MessierClassService {
	return &messierClassService{}
}
