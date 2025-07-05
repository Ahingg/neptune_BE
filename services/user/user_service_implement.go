package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	model "neptune/backend/models/user"
	"neptune/backend/pkg/jwt"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	userRepo "neptune/backend/repositories/user"
	"os"
	"strings"
	"time"
)

type userService struct {
	userRepo    userRepo.UserRepository
	labLogOnSvc log_on.LogOnService // Service to handle Lab API authentication
	labMeSvc    me.MeService        // Service to handle Lab API user profile
}

// NewUserService creates a new instance of UserService
func NewUserService(
	userRepo userRepo.UserRepository,
	labLogOnSvc log_on.LogOnService,
	labMeSvc me.MeService,
) UserService {
	return &userService{
		userRepo:    userRepo,
		labLogOnSvc: labLogOnSvc,
		labMeSvc:    labMeSvc,
	}
}

func (s *userService) LoginAssistant(ctx context.Context, req *requests.LoginRequest) (*responses.LoginResponse, string, time.Time, error) {

	// Pake Login Asssitant dlu
	logOnResp, err := s.labLogOnSvc.LogOnAssistant(ctx, req.Username, req.Password)

	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to log on assistant: %w", err)
	}

	// Di titik ini berarti udah dpt tokennya, tinggal pake /me
	meResp, err := s.labMeSvc.GetAssistantProfile(ctx, logOnResp.AccessToken)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to get assistant profile: %w", err)
	}

	// Di titik ini berarti udah berhasil dapetin data user, tinggal disesuain aja datanya trus balikin
	if s.isUserAdmin(meResp.Username) {
		meResp.Role = model.RoleAdmin
	}

	// Save data dari usernya ke db, biar ga buang buang waktu kalau nanti mau ngeseed.
	selfToken, err := jwt.CreateJWT(
		meResp.UserID,
		meResp.Username,
		meResp.Name,
		meResp.Role,
		time.Now().Add(time.Duration(logOnResp.ExpiresIn)*time.Second),
	)

	// later save messier token to the db.

	return &responses.LoginResponse{
			UserID:   meResp.UserID,
			Username: meResp.Username,
			Name:     meResp.Name,
			Role:     meResp.Role.String(),
		},
		selfToken,
		time.Now().Add(time.Duration(logOnResp.ExpiresIn) * time.Second),
		nil
}

func (s *userService) LoginStudent(ctx context.Context, req *requests.LoginRequest) (*responses.LoginResponse, string, time.Time, error) {

	// Pake Login Student dlu
	logOnResp, err := s.labLogOnSvc.LogOnStudent(ctx, req.Username, req.Password)

	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to log on student: %w", err)
	}

	selfToken, err := jwt.CreateJWT(
		logOnResp.Student.UserID.String(),
		logOnResp.Student.UserName,
		logOnResp.Student.Name,
		model.RoleStudent,
		logOnResp.Token.Expires,
	)

	// Di titik ini berarti udah berhasil dapetin data user, tinggal disesuain aja datanya trus balikin
	return &responses.LoginResponse{
			UserID:   logOnResp.Student.UserID.String(),
			Username: logOnResp.Student.UserName,
			Name:     logOnResp.Student.Name,
			Role:     model.RoleStudent.String(),
		},
		selfToken,
		logOnResp.Token.Expires,
		nil
}

// GetUserProfile retrieves user profile by ID from the internal database.
func (s *userService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID) // Assuming GetUserByID exists in repo
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile from repository: %w", err)
	}
	if u == nil {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *userService) isUserAdmin(username string) bool {
	return strings.EqualFold(username, os.Getenv("ADMIN_USERNAME"))
}
