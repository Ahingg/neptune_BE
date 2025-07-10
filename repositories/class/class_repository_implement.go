package class

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	models "neptune/backend/models/class"
)

type classRepositoryImplement struct {
	db *gorm.DB
}

func (c classRepositoryImplement) AddClassStudents(ctx context.Context, classTransactionID string, studentUserIDs []uuid.UUID) error {
	if len(studentUserIDs) == 0 {
		return nil
	}
	var classStudents []models.ClassStudent
	for _, userID := range studentUserIDs {
		id, err := uuid.Parse(classTransactionID)
		if err != nil {
			return fmt.Errorf("invalid class transaction ID %s: %w", classTransactionID, err)
		}
		classStudents = append(classStudents, models.ClassStudent{
			ClassTransactionID: id,
			UserID:             userID,
		})
	}
	return c.db.WithContext(ctx).Create(&classStudents).Error
}

func (c classRepositoryImplement) AddClassAssistants(ctx context.Context, classTransactionID string, assistantUserIDs []uuid.UUID) error {
	if len(assistantUserIDs) == 0 {
		return nil
	}
	var classAssistants []models.ClassAssistant
	for _, userID := range assistantUserIDs {
		id, err := uuid.Parse(classTransactionID)
		if err != nil {
			return fmt.Errorf("invalid class transaction ID %s: %w", classTransactionID, err)
		}
		classAssistants = append(classAssistants, models.ClassAssistant{
			ClassTransactionID: id,
			UserID:             userID,
		})
	}
	return c.db.WithContext(ctx).Create(&classAssistants).Error
}

func (c classRepositoryImplement) ClearClassStudents(ctx context.Context, classTransactionID string) error {
	return c.db.WithContext(ctx).Where("class_transaction_id = ?", classTransactionID).Delete(&models.ClassStudent{}).Error
}

func (c classRepositoryImplement) ClearClassAssistants(ctx context.Context, classTransactionID string) error {
	return c.db.WithContext(ctx).Where("class_transaction_id = ?", classTransactionID).Delete(&models.ClassAssistant{}).Error
}

func (c classRepositoryImplement) FindClassByTransactionID(ctx context.Context, classTransactionID string) (*models.Class, error) {
	var class models.Class
	result := c.db.WithContext(ctx).
		Preload("Students.User").
		Preload("Assistants.User").
		Where("class_transaction_id = ?", classTransactionID).
		First(&class)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find class by transaction ID %s: %w", classTransactionID, result.Error)
	}
	return &class, nil
}

func (c classRepositoryImplement) FindClassBasicInfoBySemesterAndCourse(ctx context.Context, semesterID, courseOutlineID string) ([]models.Class, error) {
	var classes []models.Class
	result := c.db.WithContext(ctx).
		Select("class_transaction_id", "semester_id", "course_outline_id", "class_code").
		Where("semester_id = ?", semesterID).
		Where("course_outline_id = ?", courseOutlineID).
		Find(&classes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve basic class info by semester %s and course %s: %w", semesterID, courseOutlineID, result.Error)
	}
	return classes, nil
}

func (c classRepositoryImplement) SaveClass(ctx context.Context, class *models.Class) error {
	return c.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "class_transaction_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"semester_id":       class.SemesterID,
			"course_outline_id": class.CourseOutlineID,
			"class_code":        class.ClassCode,
		}),
	}).Create(class).Error
}

func (c classRepositoryImplement) FindAllClassesBySemesterAndCourse(ctx context.Context, semesterId string, courseId string) ([]models.Class, error) {
	var classes []models.Class
	result := c.db.WithContext(ctx).
		Preload("Students.User").
		Preload("Assistants.User").
		Where("semester_id = ?", semesterId).
		Where("course_outline_id = ?", courseId).
		Find(&classes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve classes by semester %s and course %s: %w", semesterId, courseId, result.Error)
	}
	return classes, nil
}

func (c classRepositoryImplement) FindFirstStudentByClassTransactionID(ctx context.Context, classTransactionID string) (*models.ClassStudent, error) {
	var classStudent models.ClassStudent
	result := c.db.WithContext(ctx).
		Where("class_transaction_id = ?", classTransactionID).
		Preload("User"). // <-- Add the field name here
		First(&classStudent)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find first student by class transaction ID %s: %w", classTransactionID, result.Error)
	}
	return &classStudent, nil
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepositoryImplement{
		db: db,
	}
}
