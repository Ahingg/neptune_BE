package submissionServ

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	submissionRepo "neptune/backend/repositories/submission"
	"os"
	"path/filepath"
	"strings"
)

// SubmissionReviewService defines the interface for reviewing submission code.
type SubmissionReviewService interface {
	GetSubmissionCode(ctx context.Context, submissionID uuid.UUID) ([]byte, string, error)
	GetSubmissionCodeAsZip(ctx context.Context, submissionID uuid.UUID) (*bytes.Buffer, string, error)
}

type submissionReviewServiceImpl struct {
	submissionRepo submissionRepo.SubmissionRepository
}

// NewSubmissionReviewService creates a new instance of the review service.
func NewSubmissionReviewService(submissionRepo submissionRepo.SubmissionRepository) SubmissionReviewService {
	return &submissionReviewServiceImpl{submissionRepo: submissionRepo}
}

// GetSubmissionCode retrieves the raw source code and its content type.
func (s *submissionReviewServiceImpl) GetSubmissionCode(ctx context.Context, submissionID uuid.UUID) ([]byte, string, error) {
	submission, err := s.submissionRepo.FindByID(ctx, submissionID.String())
	if err != nil {
		return nil, "", fmt.Errorf("submission with ID %s not found: %w", submissionID, err)
	}

	// The path stored in DB is like "/public/submissions/..."
	// We need a local file system path, so we trim the leading slash.
	filePath := strings.TrimPrefix(submission.SourceCodePath, "/")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("source code file not found on disk for submission %s", submissionID)
	}

	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read source code file: %w", err)
	}

	// Determine content type for proper browser rendering
	contentType := "text/plain; charset=utf-8"
	if strings.HasSuffix(filePath, ".py") {
		contentType = "text/x-python; charset=utf-8"
	} else if strings.HasSuffix(filePath, ".go") {
		contentType = "text/x-go; charset=utf-8"
	}

	return code, contentType, nil
}

// GetSubmissionCodeAsZip creates a zip archive containing the submission's source code.
func (s *submissionReviewServiceImpl) GetSubmissionCodeAsZip(ctx context.Context, submissionID uuid.UUID) (*bytes.Buffer, string, error) {
	submission, err := s.submissionRepo.FindByID(ctx, submissionID.String())
	if err != nil {
		return nil, "", fmt.Errorf("submission with ID %s not found: %w", submissionID, err)
	}

	filePath := strings.TrimPrefix(submission.SourceCodePath, "/")
	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read source code file for zipping: %w", err)
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create a file within the zip archive.
	originalFilename := filepath.Base(filePath) // e.g., "main.py"
	zipFile, err := zipWriter.Create(originalFilename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create file in zip archive: %w", err)
	}

	// Write the code content to the file in the zip.
	_, err = io.Copy(zipFile, bytes.NewReader(code))
	if err != nil {
		return nil, "", fmt.Errorf("failed to write content to zip file: %w", err)
	}

	// It's important to close the writer to finalize the archive.
	zipWriter.Close()

	// Define the name of the downloaded zip file.
	downloadFilename := fmt.Sprintf("submission_%s.zip", submissionID.String())

	return buf, downloadFilename, nil
}
