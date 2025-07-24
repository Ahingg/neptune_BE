package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	model "neptune/backend/models/user"
	"neptune/backend/pkg/jwt"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	classRepository "neptune/backend/repositories/class"
	messierTokenRepo "neptune/backend/repositories/messier_token"
	"neptune/backend/repositories/semester"
	userRepo "neptune/backend/repositories/user"
	"os"
	"strings"
	"time"
)

type userService struct {
	userRepo         userRepo.UserRepository
	labLogOnSvc      log_on.LogOnService // Service to handle Lab API authentication
	labMeSvc         me.MeService        // Service to handle Lab API user profile
	messierTokenRepo messierTokenRepo.MessierTokenRepository
	classRepo        classRepository.ClassRepository
	semesterRepo     semester.SemesterRepository // Repository for Messier tokens
}

// NewUserService creates a new instance of UserService
func NewUserService(
	userRepo userRepo.UserRepository,
	labLogOnSvc log_on.LogOnService,
	labMeSvc me.MeService,
	messierTokenRepository messierTokenRepo.MessierTokenRepository,
	classRepo classRepository.ClassRepository,
	semesterRepo semester.SemesterRepository,
) UserService {
	return &userService{
		userRepo:         userRepo,
		labLogOnSvc:      labLogOnSvc,
		labMeSvc:         labMeSvc,
		messierTokenRepo: messierTokenRepository,
		classRepo:        classRepo,
		semesterRepo:     semesterRepo,
	}
}

func (s *userService) getUserEnrollmentsInCurrentSemester(ctx context.Context, userID uuid.UUID) ([]responses.UserEnrollmentDetail, error) {
	classStudents, err := s.classRepo.FindClassesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve enrollments for user %s: %w", userID.String(), err)
	}

	currentSemester, err := s.semesterRepo.FindCurrentSemester(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve current semester: %w", err)
	}

	var enrollments []responses.UserEnrollmentDetail
	for _, cs := range classStudents {
		// Ensure Class association was loaded and is not a zero value (if PK is UUID)
		// Assuming Class.ExternalClassTransactionID is a uuid.UUID, its zero value is uuid.Nil
		if cs.Class.ClassTransactionID.String() != uuid.Nil.String() && cs.Class.SemesterID.String() == currentSemester.ID {
			enrollments = append(enrollments, responses.UserEnrollmentDetail{
				ClassTransactionID: cs.Class.ClassTransactionID.String(),
				ClassName:          cs.Class.ClassCode,
				CourseOutlineID:    cs.Class.CourseOutlineID.String(), // Convert UUID to string for DTO
				SemesterID:         cs.Class.SemesterID.String(),      // Convert UUID to string for DTO
			})
		} else if cs.Class.ClassTransactionID.String() == uuid.Nil.String() {
			// This might indicate a problem with preloading or data inconsistency
			log.Printf("Warning: Class association not loaded for ClassStudent ID %d (user %s)", cs.ClassTransactionID, userID.String())
		}
	}
	return enrollments, nil
}

func (s *userService) LoginAssistant(ctx context.Context, req *requests.LoginRequest) (*responses.LoginResponse, string, time.Time, error) {
	var (
		messierAccessToken  string
		messierTokenExpires time.Time
		meResp              *me.MeResponse
		err                 error
	)

	// 1. Try to find existing MessierToken in DB
	// We need the internal UserID to look up the MessierToken.
	// First, try to find the user by username (NIM).
	internalUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to look up internal user by username: %w", err)
	}

	if internalUser != nil {
		// User exists internally, try to use their cached Messier token
		messierToken, err := s.messierTokenRepo.GetMessierTokenByUserID(ctx, internalUser.ID.String())
		if err != nil {
			log.Printf("Warning: Failed to retrieve cached Messier token for user %s: %v", internalUser.ID.String(), err)
			// Proceed to get new token if DB lookup fails
		} else if messierToken != nil && messierToken.MessierTokenExpires.After(time.Now().Add(5*time.Minute)) {
			// Token found and is not expired (add a buffer, e.g., 5 minutes)
			messierAccessToken = messierToken.MessierAccessToken
			messierTokenExpires = messierToken.MessierTokenExpires
			log.Printf("Using cached Messier token for user %s", internalUser.Username)

			// Try to get profile with cached token
			meResp, err = s.labMeSvc.GetAssistantProfile(ctx, messierAccessToken)
			if err != nil {
				log.Printf("Cached Messier token failed for user %s: %v. Attempting new login.", internalUser.Username, err)
				// If profile fetch fails, token might be invalid/revoked, proceed to full login
				messierAccessToken = "" // Clear token to force new login
			}
		}
	}

	// 2. If no valid cached token or profile fetch failed, perform full login
	if messierAccessToken == "" {
		logOnResp, logonErr := s.labLogOnSvc.LogOnAssistant(ctx, req.Username, req.Password)
		if logonErr != nil {
			return nil, "", time.Time{}, fmt.Errorf("failed to log on assistant: %w", logonErr)
		}
		messierAccessToken = logOnResp.AccessToken
		messierTokenExpires = time.Now().Add(time.Duration(logOnResp.ExpiresIn) * time.Second)

		meResp, err = s.labMeSvc.GetAssistantProfile(ctx, messierAccessToken)
		if err != nil {
			return nil, "", time.Time{}, fmt.Errorf("failed to get assistant profile with new token: %w", err)
		}
		log.Printf("Successfully obtained new Messier token for user %s", req.Username)
	}

	// 3. Determine user's role (from Messier or internal logic)
	userRole := model.RoleAssistant // Default role
	if s.isUserAdmin(meResp.Username) {
		userRole = model.RoleAdmin
	} else if meResp.Role != "" { // Use role from Messier if available and not empty
		userRole = model.Role(meResp.Role) // Convert string to model.Role
	}

	// 4. Upsert/Update User record in our DB
	if internalUser == nil {
		// User does not exist internally, create new
		internalUser = &model.User{
			ID:        uuid.New(), // Generate new UUID for internal user
			Username:  meResp.Username,
			Name:      meResp.Name,
			Role:      userRole,
			CreatedAt: time.Now(),
		}
		log.Printf("Creating new internal user record for %s", internalUser.Username)
	} else {
		// User exists internally, update details if changed
		if internalUser.Name != meResp.Name || internalUser.Role != userRole {
			internalUser.Name = meResp.Name
			internalUser.Role = userRole
			log.Printf("Updating existing internal user record for %s", internalUser.Username)
		} else {
			log.Printf("Internal user record for %s is up-to-date.", internalUser.Username)
		}
	}

	if err := s.userRepo.Save(ctx, internalUser); err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to save/update internal user record: %w", err)
	}

	// 5. Save/Update Messier token in our DB (linked to internal UserID)
	messierTokenRecord := &model.MessierToken{
		UserID:              internalUser.ID.String(),
		MessierAccessToken:  messierAccessToken,
		MessierTokenExpires: messierTokenExpires,
	}
	if err := s.messierTokenRepo.Save(ctx, messierTokenRecord); err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to save messier token to DB: %w", err)
	}

	enrollments, err := s.getUserEnrollmentsInCurrentSemester(ctx, internalUser.ID)
	if err != nil {
		log.Printf("Warning: Could not fetch enrollments for user %s: %v", internalUser.Username, err)
		// Don't fail login, but log the issue if enrollments can't be fetched
	}

	// 6. Create and return our internal JWT
	selfToken, err := jwt.CreateJWT(
		internalUser.ID.String(), // Use internal UserID for our JWT
		internalUser.Username,
		internalUser.Name,
		internalUser.Role,
		time.Now().Add(time.Hour*24), // Our JWT expiry, e.g., 24 hours
	)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to create internal JWT: %w", err)
	}

	return &responses.LoginResponse{
		UserID:      internalUser.ID.String(),
		Username:    internalUser.Username,
		Name:        internalUser.Name,
		Role:        internalUser.Role.String(),
		Enrollments: enrollments,
	}, selfToken, time.Now().Add(time.Hour * 24), nil
}

func (s *userService) LoginStudent(ctx context.Context, req *requests.LoginRequest) (*responses.LoginResponse, string, time.Time, error) {
	var (
		messierAccessToken  string
		messierTokenExpires time.Time
		logOnResp           *log_on.LogOnStudentResponse
		err                 error
	)

	// 1. Try to find existing MessierToken in DB
	internalUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to look up internal user by username: %w", err)
	}

	logOnResp, err = s.labLogOnSvc.LogOnStudent(ctx, req.Username, req.Password)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to log on student: %w", err)
	}
	messierAccessToken = logOnResp.Token.Token
	messierTokenExpires = logOnResp.Token.Expires
	log.Printf("Successfully obtained new Messier token for student %s", req.Username)

	// 3. Upsert/Update User record in our DB
	userRole := model.RoleStudent // Students always have RoleStudent
	if internalUser == nil {
		internalUser = &model.User{
			ID:        uuid.New(),
			Username:  logOnResp.Student.UserName,
			Name:      logOnResp.Student.Name,
			Role:      userRole,
			CreatedAt: time.Now(),
		}
		log.Printf("Creating new internal student record for %s", internalUser.Username)
	} else {
		if internalUser.Name != logOnResp.Student.Name || internalUser.Role != userRole {
			internalUser.Name = logOnResp.Student.Name
			internalUser.Role = userRole
			log.Printf("Updating existing internal student record for %s", internalUser.Username)
		} else {
			log.Printf("Internal student record for %s is up-to-date.", internalUser.Username)
		}
	}

	if err := s.userRepo.Save(ctx, internalUser); err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to save/update internal user record: %w", err)
	}

	// 4. Save/Update Messier token in our DB
	messierTokenRecord := &model.MessierToken{
		UserID:              internalUser.ID.String(),
		MessierAccessToken:  messierAccessToken,
		MessierTokenExpires: messierTokenExpires,
	}
	if err := s.messierTokenRepo.Save(ctx, messierTokenRecord); err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to save messier token to DB: %w", err)
	}

	enrollments, err := s.getUserEnrollmentsInCurrentSemester(ctx, internalUser.ID)
	if err != nil {
		log.Printf("Warning: Could not fetch enrollments for user %s: %v", internalUser.Username, err)
	}
	// 5. Create and return our internal JWT
	selfToken, err := jwt.CreateJWT(
		internalUser.ID.String(),
		internalUser.Username,
		internalUser.Name,
		internalUser.Role,
		time.Now().Add(time.Hour*24), // Our JWT expiry, e.g., 24 hours
	)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("failed to create internal JWT: %w", err)
	}

	return &responses.LoginResponse{
		UserID:      internalUser.ID.String(),
		Username:    internalUser.Username,
		Name:        internalUser.Name,
		Role:        internalUser.Role.String(),
		Enrollments: enrollments,
	}, selfToken, time.Now().Add(time.Hour * 24), nil
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

func (s *userService) GetDetailedUserProfile(ctx context.Context, userID uuid.UUID) (*responses.UserMeResponse, error) {
	internalUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID %s: %w", userID.String(), err)
	}
	if internalUser == nil {
		return nil, nil // User not found
	}

	// Fetch enrollments for the user
	enrollments, err := s.getUserEnrollmentsInCurrentSemester(ctx, userID)
	if err != nil {
		log.Printf("Warning: Could not fetch enrollments for user %s during MeHandler: %v", internalUser.Username, err)
		// Don't fail the MeHandler, but log the issue
	}

	return &responses.UserMeResponse{
		ID:          internalUser.ID.String(),
		Username:    internalUser.Username,
		Name:        internalUser.Name,
		Role:        internalUser.Role.String(),
		Enrollments: enrollments,
	}, nil
}

func (s *userService) DeleteUserAccessToken(ctx context.Context, userID string) error {
	// Delete the token from the database
	err := s.messierTokenRepo.DeleteByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user access token: %w", err)
	}
	return nil
}

func (s *userService) isUserAdmin(username string) bool {
	return strings.EqualFold(username, os.Getenv("ADMIN_USERNAME"))
}
