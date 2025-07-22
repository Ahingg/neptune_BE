package contestRepository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	contestModel "neptune/backend/models/contest"

	"time"
)

type contestRepositoryImpl struct {
	db *gorm.DB
}

func NewContestRepository(db *gorm.DB) ContestRepository {
	return &contestRepositoryImpl{db: db}
}

//func (r *contestRepositoryImpl) WithTransaction(f func(tx *gorm.DB) error) error {
//	return r.db.Transaction(f)
//}

// SaveContest creates or updates a Contest.
func (r *contestRepositoryImpl) SaveContest(ctx context.Context, contest *contestModel.Contest) error {
	if contest.ID == uuid.Nil {
		contest.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // Conflict on primary key (ID)
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":        contest.Name,
			"scope":       contest.Scope,
			"description": contest.Description,
			"updated_at":  time.Now(),
		}),
	}).Create(contest).Error
}

// FindContestByID retrieves a Contest with its associated Cases.
func (r *contestRepositoryImpl) FindContestByID(ctx context.Context, contestID uuid.UUID) (*contestModel.Contest, error) {
	var contest contestModel.Contest
	result := r.db.WithContext(ctx).
		Preload("ContestCases.Case"). // Preload join table, then the Case itself
		Where("id = ?", contestID).
		First(&contest)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find contest by ID %s: %w", contestID.String(), result.Error)
	}
	return &contest, nil
}

// FindAllContests retrieves all Contests (basic info).
func (r *contestRepositoryImpl) FindAllContests(ctx context.Context) ([]contestModel.Contest, error) {
	var contests []contestModel.Contest
	result := r.db.WithContext(ctx).Find(&contests)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all contests: %w", result.Error)
	}
	return contests, nil
}

// DeleteContest soft deletes a contest.
func (r *contestRepositoryImpl) DeleteContest(ctx context.Context, contestID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&contestModel.Contest{}, contestID).Error
}

// AddCasesToContest adds multiple cases to a contest (via ContestCase join table).
// It clears existing assignments for the given contest before adding new ones.
func (r *contestRepositoryImpl) AddCasesToContest(ctx context.Context, contestID uuid.UUID, contestCases []contestModel.ContestCase) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(contestCases) == 0 {
			return nil
		}

		// Fetch existing (contest_id, case_id) pairs
		var existing []contestModel.ContestCase
		if err := tx.
			Where("contest_id = ? AND case_id IN ?", contestID, getCaseIDs(contestCases)).
			Find(&existing).Error; err != nil {
			return fmt.Errorf("failed to check existing cases: %w", err)
		}

		existingMap := make(map[uuid.UUID]bool)
		for _, ec := range existing {
			existingMap[ec.CaseID] = true
		}

		// Filter out duplicates
		var filtered []contestModel.ContestCase
		for _, c := range contestCases {
			if !existingMap[c.CaseID] {
				c.ContestID = contestID
				filtered = append(filtered, c)
			}
		}

		if len(filtered) > 0 {
			if err := tx.Create(&filtered).Error; err != nil {
				return fmt.Errorf("failed to add filtered cases: %w", err)
			}
		}

		return nil
	})
}

// helper to extract []uuid.UUID from []ContestCase
func getCaseIDs(contestCases []contestModel.ContestCase) []uuid.UUID {
	caseIDs := make([]uuid.UUID, 0, len(contestCases))
	for _, c := range contestCases {
		caseIDs = append(caseIDs, c.CaseID)
	}
	return caseIDs
}

// ClearContestCases removes all cases from a contest.
func (r *contestRepositoryImpl) ClearContestCases(ctx context.Context, contestID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("contest_id = ?", contestID).Delete(&contestModel.ContestCase{}).Error
}

// AssignContestToClass assigns a contest to a class with specific start/end times.
// It upserts the ClassContest record.
func (r *contestRepositoryImpl) AssignContestToClass(ctx context.Context, classContest *contestModel.ClassContest) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "class_transaction_id"},
			{Name: "contest_id"},
		}, // Conflict on composite primary key
		DoUpdates: clause.Assignments(map[string]interface{}{
			"start_time": classContest.StartTime,
			"end_time":   classContest.EndTime,
			"updated_at": time.Now(),
		}),
	}).Create(classContest).Error
}

// FindContestsByClassTransactionID retrieves all contests assigned to a specific class, with their durations.
func (r *contestRepositoryImpl) FindContestsByClassTransactionID(ctx context.Context, classTransactionID uuid.UUID) ([]contestModel.ClassContest, error) {
	var classContests []contestModel.ClassContest
	result := r.db.WithContext(ctx).
		Preload("Contest").                   // Preload the Contest details
		Preload("Contest.ContestCases.Case"). // Further preload Cases within the Contest
		Where("class_transaction_id = ?", classTransactionID).
		Find(&classContests)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find contests for class %s: %w", classTransactionID.String(), result.Error)
	}
	return classContests, nil
}

// FindClassContestByIDs finds a specific ClassContest entry.
func (r *contestRepositoryImpl) FindClassContestByIDs(ctx context.Context, classTransactionID, contestID uuid.UUID) (*contestModel.ClassContest, error) {
	var classContest contestModel.ClassContest
	result := r.db.WithContext(ctx).
		Where("class_transaction_id = ?", classTransactionID).
		Where("contest_id = ?", contestID).
		First(&classContest)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find class contest for class %s and contest %s: %w", classTransactionID.String(), contestID.String(), result.Error)
	}
	return &classContest, nil
}

// FindContestCases retrieves all cases for a specific contest.
func (r *contestRepositoryImpl) FindContestCases(ctx context.Context, contestID uuid.UUID) ([]contestModel.ContestCase, error) {
	var cases []contestModel.ContestCase
	result := r.db.WithContext(ctx).
		Preload("Case"). // Preload the Case details
		Where("contest_id = ?", contestID).
		Find(&cases)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cases for contest %s: %w", contestID.String(), result.Error)
	}
	return cases, nil
}

func (r *contestRepositoryImpl) GetCaseCountInContest(ctx context.Context, contestID uuid.UUID) (int, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&contestModel.ContestCase{}).
		Where("contest_id = ?", contestID).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count cases in contest %s: %w", contestID.String(), result.Error)
	}
	return int(count), nil
}

// RemoveContestFromClass deletes a contest assignment from a class.
func (r *contestRepositoryImpl) RemoveContestFromClass(ctx context.Context, classTransactionID, contestID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("class_transaction_id = ?", classTransactionID).
		Where("contest_id = ?", contestID).
		Delete(&contestModel.ClassContest{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove contest %s from class %s: %w", contestID.String(), classTransactionID.String(), result.Error)
	}
	return nil
}
