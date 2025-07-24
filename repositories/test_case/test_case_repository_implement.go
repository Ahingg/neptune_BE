package testCaseRepo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	testCaseModel "neptune/backend/models/test_case"
	"time"
)

type testCaseRepository struct {
	db *gorm.DB
}

func (t *testCaseRepository) SaveTestCase(ctx context.Context, testCase *testCaseModel.TestCase) error {
	return t.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "case_id"}, {Name: "number"}}, // Conflict on composite primary key
		DoUpdates: clause.Assignments(map[string]interface{}{
			"input_url":  testCase.InputUrl,
			"output_url": testCase.OutputUrl,
			"created_at": time.Now(), // Update creation time to reflect latest upload
		}),
	}).Create(testCase).Error
}

func (t *testCaseRepository) SaveTestCaseBatch(ctx context.Context, testCases []testCaseModel.TestCase) error {
	if len(testCases) == 0 {
		return nil
	}
	return t.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "case_id"}, {Name: "number"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"input_url":  gorm.Expr("EXCLUDED.input_url"),
			"output_url": gorm.Expr("EXCLUDED.output_url"),
			"created_at": gorm.Expr("EXCLUDED.created_at"),
		}),
	}).CreateInBatches(testCases, 100).Error // Batch size 100
}

func (t *testCaseRepository) FindTestCaseByCaseID(ctx context.Context, caseID string) ([]testCaseModel.TestCase, error) {
	var testcases []testCaseModel.TestCase
	result := t.db.WithContext(ctx).Where("case_id = ?", caseID).Order("number").Find(&testcases)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find testcases for case ID %s: %w", caseID, result.Error)
	}
	return testcases, nil
}

func (t *testCaseRepository) DeleteTestCaseByCaseID(ctx context.Context, caseID string) error {
	return t.db.WithContext(ctx).Unscoped().Where("case_id = ?", caseID).Delete(&testCaseModel.TestCase{}).Error
}

// NewTestCaseRepository creates a new instance of TestCaseRepository
func NewTestCaseRepository(db *gorm.DB) TestCaseRepository {
	return &testCaseRepository{
		db: db,
	}
}
