package internal_class

import (
	"context"
	"fmt"
	"log"
	messierClass "neptune/backend/messier/class"
	"neptune/backend/messier/constants"
	models "neptune/backend/models/class"
	model "neptune/backend/models/user"
	"neptune/backend/pkg/responses"
	"neptune/backend/pkg/utils"
	classRepository "neptune/backend/repositories/class"
	"neptune/backend/repositories/messier_token"
	userRepository "neptune/backend/repositories/user"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type classService struct {
	messierClassSrv  messierClass.MessierClassService
	classRepo        classRepository.ClassRepository
	userRepo         userRepository.UserRepository
	messierTokenRepo messier_token.MessierTokenRepository
}

func (c classService) SyncClasses(ctx context.Context, semesterID string, requestMakerID string) error {
	authToken, err := utils.GetAndValidateMessierToken(ctx, requestMakerID, c.messierTokenRepo)
	if err != nil {
		return err
	}

	courseOutlineIDs := []string{
		constants.AlgoprogID1,
		constants.AlgoprogID2,
	}
	totalSyncedClasses := 0
	for _, courseID := range courseOutlineIDs {
		log.Printf("Starting basic class sync for Semester: %s, CourseOutline: %s", semesterID, courseID)

		basicClasses, err := c.messierClassSrv.GetClassesBySemesterAndCourseOutline(ctx, semesterID, courseID, authToken)
		if err != nil {
			log.Printf("Error fetching basic classes for semester %s, course %s: %v", semesterID, courseID, err)
			continue
		}

		for _, bc := range basicClasses {
			semId, err := uuid.Parse(semesterID)
			if err != nil {
				log.Printf("Invalid SemesterID on sync class: %s, CourseOutlineID: %s, ClassCode: %s", semesterID, bc.CourseOutlineID, bc.ClassCode)
				continue
			}
			class := &models.Class{
				ClassTransactionID: bc.ClassTransactionID,
				SemesterID:         semId,
				CourseOutlineID:    bc.CourseOutlineID,
				ClassCode:          bc.ClassCode,
			}
			if err := c.classRepo.SaveClass(ctx, class); err != nil {
				log.Printf("Error saving basic class %s (%s): %v", bc.ClassTransactionID, bc.ClassCode, err)
			} else {
				log.Printf("Synced basic class: %s - %s", bc.ClassTransactionID, bc.ClassCode)
				totalSyncedClasses++
			}
		}
	}
	log.Printf("Successfully synced %d basic classes for semester %s.", totalSyncedClasses, semesterID)
	return nil
}

func (c classService) SyncClassStudents(ctx context.Context, semesterID string, courseOutlineID string, requestMakerID string) error {
	messierAccessToken, err := utils.GetAndValidateMessierToken(ctx, requestMakerID, c.messierTokenRepo)
	if err != nil {
		return err
	}

	// 1. Get all basic classes from internal DB that need student syncing
	classes, err := c.classRepo.FindClassBasicInfoBySemesterAndCourse(ctx, semesterID, courseOutlineID)
	if err != nil {
		return fmt.Errorf("failed to retrieve basic classes for student sync: %w", err)
	}

	totalSyncedStudents := 0
	for _, cl := range classes {
		log.Printf("Syncing students for Class: %s (Name: %s)", cl.ClassTransactionID, cl.ClassCode)

		// Fetch students from Messier
		messierStudents, err := c.messierClassSrv.GetStudentFromClassTransaction(ctx, cl.SemesterID.String(), cl.CourseOutlineID.String(), cl.ClassCode, messierAccessToken)
		if err != nil {
			log.Printf("Warning: Failed to fetch students for class %s: %v", cl.ClassTransactionID, err)
			continue
		}

		var studentUserIDs []uuid.UUID
		for _, ms := range messierStudents {
			user, err := c.userRepo.GetUserByUsername(ctx, ms.NIM)
			if err != nil {
				log.Printf("Error finding user by NIM %s: %v", ms.NIM, err)
				continue
			}

			if user == nil {
				newUser := &model.User{
					ID:        ms.BinusianID,
					Username:  ms.NIM,
					Name:      ms.Name,
					Role:      model.RoleStudent,
					CreatedAt: time.Now(),
				}
				if err := c.userRepo.CreateUser(ctx, newUser); err != nil {
					log.Printf("Error creating new student user %s (%s): %v", ms.NIM, ms.Name, err)
					continue
				}
				user = newUser
				log.Printf("Created new student user: %s (%s)", user.Username, user.Name)
			} else {
				// Update existing user details if necessary
				if user.Name != ms.Name || user.ID != ms.BinusianID || user.Role != model.RoleStudent {
					user.Name = ms.Name
					user.ID = ms.BinusianID
					user.Role = model.RoleStudent
					if err := c.userRepo.UpdateUser(ctx, user); err != nil {
						log.Printf("Warning: Failed to update existing student user %s (%s): %v", user.Username, user.Name, err)
					}
				}
			}
			studentUserIDs = append(studentUserIDs, user.ID)
		}

		// Use a transaction for clearing and adding students for atomicity per class
		err = c.classRepo.ClearClassStudents(ctx, cl.ClassTransactionID.String())
		if err != nil {
			log.Printf("Error clearing existing students for class %s: %v", cl.ClassTransactionID, err)
			continue
		}
		err = c.classRepo.AddClassStudents(ctx, cl.ClassTransactionID.String(), studentUserIDs)
		if err != nil {
			log.Printf("Error adding new students for class %s: %v", cl.ClassTransactionID.String(), err)
			continue
		}
		log.Printf("Synced %d students for class %s.", len(studentUserIDs), cl.ClassTransactionID.String())
		totalSyncedStudents += len(studentUserIDs)
	}
	log.Printf("Successfully synced students for %d classes. Total students synced: %d", len(classes), totalSyncedStudents)
	return nil
}

// SyncClassAssistants : Please make sure to sync the student first
func (c classService) SyncClassAssistants(ctx context.Context, semesterID string, courseOutlineID, requestMakerID string) error {
	authToken, err := utils.GetAndValidateMessierToken(ctx, requestMakerID, c.messierTokenRepo)
	if err != nil {
		return fmt.Errorf("failed to get and validate Messier token: %w", err)
	}

	classes, err := c.classRepo.FindClassBasicInfoBySemesterAndCourse(ctx, semesterID, courseOutlineID)
	if err != nil {
		return fmt.Errorf("failed to retrieve basic classes for assistant sync: %w", err)
	}

	generationRegex := regexp.MustCompile(`[A-Z]{2}(\d{2}-\d{1})`)

	totalSyncedAssitants := 0
	for _, cl := range classes {
		log.Printf("Syncing assistants for Class: %s (Name: %s)", cl.ClassTransactionID, cl.ClassCode)

		// Fetch initial assistants from Messier
		firstStudent, err := c.classRepo.FindFirstStudentByClassTransactionID(ctx, cl.ClassTransactionID.String())
		if err != nil {
			log.Printf("Error finding first student for class %s: %v", cl.ClassTransactionID, err)
			continue
		}

		assistants, err := c.messierClassSrv.GetAssistantInitialFromStudentTransaction(ctx, firstStudent.User.Username, semesterID, authToken)
		if err != nil {
			log.Printf("Warning: Failed on fetch Assistant initial on Username: %s and class: %s: %v", firstStudent.User.Username, cl.ClassCode, err)
			continue
		}

		var classAssistantInitials []string
		for _, assistant := range assistants {
			if assistant.ClassTransactionID == cl.ClassTransactionID {
				fmt.Printf("Found assistant %s for class %s\n", assistant.Assistants, cl.ClassTransactionID)
				classAssistantInitials = assistant.Assistants
				break
			}
		}

		if len(classAssistantInitials) == 0 {
			log.Printf("No assistants found for class %s", cl.ClassTransactionID)
			if err := c.classRepo.ClearClassAssistants(ctx, cl.ClassTransactionID.String()); err != nil {
				log.Printf("Error clearing existing assistants for class %s: %v", cl.ClassTransactionID.String(), err)
			}
			continue
		}

		var assistantUserIDs []uuid.UUID
		log.Printf("Processing %d assistant initials for class %s: %v", len(classAssistantInitials), cl.ClassCode, classAssistantInitials)
		for _, initial := range classAssistantInitials {
			matches := generationRegex.FindStringSubmatch(initial)
			if len(matches) < 2 {
				log.Printf("Warning: Could not extract generation from assistant initial %s for class %s. Skipping.", initial, cl.ClassTransactionID)
				continue
			}
			generation := matches[1]
			log.Printf("Processing assistant initial %s with generation %s for class %s", initial, generation, cl.ClassCode)

			// 2. Fetch detailed assistant info
			assistantDetail, err := c.messierClassSrv.GetAssistantDetailFromAssistantInitial(ctx, initial, generation, authToken)
			if err != nil {
				log.Printf("Warning: Failed to fetch detail for assistant %s (%s): %v", initial, generation, err)
				continue
			}

			user, err := c.userRepo.GetUserByUsername(ctx, assistantDetail.Username)
			if err != nil {
				log.Printf("Error finding user by Username %s: %v", assistantDetail.Username, err)
				continue
			}

			if user == nil {
				newUser := &model.User{
					ID:        assistantDetail.UserID,
					Username:  assistantDetail.Username,
					Name:      assistantDetail.Name,
					Role:      model.RoleAssistant,
					CreatedAt: time.Now(),
				}
				if err := c.userRepo.CreateUser(ctx, newUser); err != nil {
					log.Printf("Error creating new assistant user %s (%s): %v", assistantDetail.Username, assistantDetail.Name, err)
					continue
				}
				user = newUser
				log.Printf("Created new assistant user: %s (%s)", user.Username, user.Name)
			} else {
				// Update existing user details if necessary
				if user.Name != assistantDetail.Name || user.ID != assistantDetail.UserID || user.Role != model.RoleAssistant {
					user.Name = assistantDetail.Name
					user.ID = assistantDetail.UserID
					user.Role = model.RoleAssistant
					if err := c.userRepo.CreateUser(ctx, user); err != nil {
						log.Printf("Warning: Failed to update existing assistant user %s (%s): %v", user.Username, user.Name, err)
					}
				}
			}
			assistantUserIDs = append(assistantUserIDs, user.ID)
		}
		err = c.classRepo.ClearClassAssistants(ctx, cl.ClassTransactionID.String())
		if err != nil {
			log.Printf("Error clearing existing assistants for class %s: %v", cl.ClassTransactionID, err)
			continue
		}
		err = c.classRepo.AddClassAssistants(ctx, cl.ClassTransactionID.String(), assistantUserIDs)
		if err != nil {
			log.Printf("Error adding new assistants for class %s: %v", cl.ClassTransactionID.String(), err)
		}

		log.Printf("Synced %d assistants for class %s.", len(assistantUserIDs), cl.ClassTransactionID.String())
		totalSyncedAssitants += len(assistantUserIDs)
	}
	log.Printf("Successfully synced assistants for %d classes. Total assistants synced: %d", len(classes), totalSyncedAssitants)
	return nil
}

func (c classService) GetClassesBySemesterAndCourse(ctx context.Context, semesterID string, courseID string) ([]responses.GetClassWithoutDetailResponse, error) {
	classes, err := c.classRepo.FindAllClassesBySemesterAndCourse(ctx, semesterID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve classes for semester %s and course %s: %w", semesterID, courseID, err)
	}

	var response []responses.GetClassWithoutDetailResponse
	for _, class := range classes {
		response = append(response, responses.GetClassWithoutDetailResponse{
			ClassTransactionID: class.ClassTransactionID.String(),
			SemesterID:         class.SemesterID.String(),
			CourseOutlineID:    class.CourseOutlineID.String(),
			ClassCode:          class.ClassCode,
		})
	}
	return response, nil
}

func (c classService) GetClassDetailBySemesterAndCourse(ctx context.Context, semesterID string, courseID string) ([]responses.GetDetailClassResponse, error) {
	classes, err := c.classRepo.FindAllClassesBySemesterAndCourse(ctx, semesterID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve classes for semester %s and course %s: %w", semesterID, courseID, err)
	}

	var response []responses.GetDetailClassResponse
	for _, class := range classes {
		response = append(response, responses.GetDetailClassResponse{
			ClassTransactionID: class.ClassTransactionID.String(),
			SemesterID:         class.SemesterID.String(),
			CourseOutlineID:    class.CourseOutlineID.String(),
			ClassCode:          class.ClassCode,
			Students:           class.Students,
			Assistants:         class.Assistants,
		})
	}
	return response, nil
}

func (c classService) GetClassDetailBySemesterCourseAndStudent(ctx context.Context, semesterID, courseID, userID string) ([]responses.GetDetailClassResponse, error) {
	classes, err := c.classRepo.FindClassBySemesterCourseAndStudent(ctx, semesterID, courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve classes for semester %s, course %s, and user %s: %w", semesterID, courseID, userID, err)
	}

	var response []responses.GetDetailClassResponse
	for _, class := range classes {
		response = append(response, responses.GetDetailClassResponse{
			ClassTransactionID: class.ClassTransactionID.String(),
			SemesterID:         class.SemesterID.String(),
			CourseOutlineID:    class.CourseOutlineID.String(),
			ClassCode:          class.ClassCode,
			Students:           class.Students,
			Assistants:         class.Assistants,
		})
	}
	return response, nil
}

func (c classService) GetClassDetailByTransactionID(ctx context.Context, classTransactionID string) (*responses.GetDetailClassResponse, error) {
	class, err := c.classRepo.FindClassByTransactionID(ctx, classTransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find class by transaction ID %s: %w", classTransactionID, err)
	}
	if class == nil {
		return nil, fmt.Errorf("class with transaction ID %s not found", classTransactionID)
	}

	response := &responses.GetDetailClassResponse{
		ClassTransactionID: class.ClassTransactionID.String(),
		SemesterID:         class.SemesterID.String(),
		CourseOutlineID:    class.CourseOutlineID.String(),
		ClassCode:          class.ClassCode,
		Students:           class.Students,
		Assistants:         class.Assistants,
	}
	return response, nil
}

func NewClassService(
	messierClassSrv messierClass.MessierClassService,
	classRepo classRepository.ClassRepository,
	userRepo userRepository.UserRepository,
	messierTokenRepo messier_token.MessierTokenRepository,
) ClassService {
	return &classService{
		messierClassSrv:  messierClassSrv,
		classRepo:        classRepo,
		userRepo:         userRepo,
		messierTokenRepo: messierTokenRepo,
	}
}
