package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

var languageIdToExtensions = map[int][]string{
	76: {".cpp", ".c"}, // C++ (Clang 7.0.1)
	52: {".cpp", ".c"}, // C++ (GCC 7.4.0)
	53: {".cpp", ".c"}, // C++ (GCC 8.3.0)
	54: {".cpp", ".c"}, // C++ (GCC 9.2.0)
	75: {".cpp", ".c"}, // C (Clang 7.0.1)
	48: {".cpp", ".c"}, // C (GCC 7.4.0)
	49: {".cpp", ".c"}, // C (GCC 8.3.0)
	50: {".cpp", ".c"}, // C (GCC 9.2.0)
	92: {".py"},        // Python
}

type SubmitCodeRequest struct {
	CaseID     uuid.UUID
	LanguageID int
	ContestID  *uuid.UUID

	// Internally populated fields after parsing
	SourceCodeBytes []byte
	FileExtension   string
}

// ParseAndValidate handles the logic of parsing a multipart/form-data request.
// It populates the request struct and performs validation, keeping the handler clean.
func (r *SubmitCodeRequest) ParseAndValidate(c *gin.Context) error {
	// --- Field Parsing ---
	caseIDStr := c.PostForm("case_id")
	langIDStr := c.PostForm("language_id")
	contestIDStr := c.PostForm("contest_id")
	sourceCodeStr := c.PostForm("source_code")

	if caseIDStr == "" || langIDStr == "" {
		return fmt.Errorf("case_id and language_id are required fields")
	}

	var err error
	if r.CaseID, err = uuid.Parse(caseIDStr); err != nil {
		return fmt.Errorf("invalid case_id format: %w", err)
	}
	if r.LanguageID, err = strconv.Atoi(langIDStr); err != nil {
		return fmt.Errorf("invalid language_id format: %w", err)
	}
	if contestIDStr != "" {
		parsedContestID, err := uuid.Parse(contestIDStr)
		if err != nil {
			return fmt.Errorf("invalid contest_id format: %w", err)
		}
		r.ContestID = &parsedContestID
	}

	// --- Source Code Handling (File or String) ---
	file, err := c.FormFile("source_file")
	isPostWithFile := err == nil

	if sourceCodeStr != "" && isPostWithFile {
		return fmt.Errorf("provide either source_code string or source_file, not both")
	}
	if sourceCodeStr == "" && !isPostWithFile {
		return fmt.Errorf("either source_code string or source_file must be provided")
	}

	// Validate extension and read content
	acceptedExtensions, ok := languageIdToExtensions[r.LanguageID]
	if !ok {
		return fmt.Errorf("language_id %d is not supported", r.LanguageID)
	}

	if isPostWithFile {
		// --- File Logic ---
		actualExtension := strings.ToLower(filepath.Ext(file.Filename))
		isValidExt := false
		for _, ext := range acceptedExtensions {
			if actualExtension == ext {
				isValidExt = true
				break
			}
		}
		if !isValidExt {
			return fmt.Errorf("file extension mismatch: expected one of %v for the selected language, but got '%s'", acceptedExtensions, actualExtension)
		}

		r.FileExtension = actualExtension

		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open uploaded file: %w", err)
		}
		defer srcFile.Close()

		r.SourceCodeBytes, err = io.ReadAll(srcFile)
		if err != nil {
			return fmt.Errorf("failed to read uploaded file: %w", err)
		}
	} else {
		// --- String Logic ---
		r.SourceCodeBytes = []byte(sourceCodeStr)
		// Assign default extension
		r.FileExtension = acceptedExtensions[0]
	}

	return nil
}
