package semester

import "context"

type MessierSemesterService interface {
	GetSemesters(ctx context.Context, authToken string) ([]GetSemestersResponse, error)
}
