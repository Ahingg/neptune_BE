package log_on

import (
	"context"
)

type LogOnService interface {
	LogOnStudent(ctx context.Context, username, password string) (response *LogOnStudentResponse, err error)
	LogOnAssistant(ctx context.Context, username, password string) (response *LogOnResponse, err error)
}
