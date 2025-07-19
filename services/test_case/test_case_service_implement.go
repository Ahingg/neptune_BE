package testCaseServ

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	testCaseModel "neptune/backend/models/test_case"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	caseRepository "neptune/backend/repositories/case"
	testCaseRepo "neptune/backend/repositories/test_case"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testcaseServiceImpl struct {
	testcaseRepo testCaseRepo.TestCaseRepository
	caseRepo     caseRepository.CaseRepository
}

func (s testcaseServiceImpl) UploadTestCases(ctx context.Context, req requests.AddTestCaseRequest) error {
	// 1. Verify the Case exists
	problemCase, err := s.caseRepo.FindCaseByID(ctx, req.CaseID)
	if err != nil {
		return fmt.Errorf("failed to find parent case %s: %w", req.CaseID.String(), err)
	}
	if problemCase == nil {
		return fmt.Errorf("case with ID %s not found", req.CaseID.String())
	}

	// 2. Clear existing testcases for this problem
	if err := s.testcaseRepo.DeleteTestCaseByCaseID(ctx, req.CaseID.String()); err != nil {
		return fmt.Errorf("failed to clear existing testcases for case %s: %w", req.CaseID.String(), err)
	}

	// 3. Open the uploaded zip file
	src, err := req.File.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded zip file: %w", err)
	}
	defer src.Close()

	zipReader, err := zip.NewReader(src, req.File.Size)
	if err != nil {
		return fmt.Errorf("failed to read zip file from multipart header: %w", err)
	}

	caseTestcasesDir := filepath.Join("./private/test_case", req.CaseID.String())
	if err := os.MkdirAll(caseTestcasesDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create testcase directory %s: %w", caseTestcasesDir, err)
	}
	tempTestcasesByDir := make(map[string]*struct {
		InputFile  *zip.File
		OutputFile *zip.File
	})

	// No longer need currentTestcaseNumber here; it will be assigned sequentially *after* grouping.

	for _, f := range zipReader.File {
		// Skip directories themselves, we only care about files.
		if f.FileInfo().IsDir() {
			continue
		}

		// Normalize the path to use forward slashes for consistent splitting
		normalizedPath := strings.ReplaceAll(f.Name, "\\", "/")

		// Extract filename and directory path
		dir, filename := filepath.Split(normalizedPath)
		dir = strings.TrimSuffix(dir, "/") // Remove trailing slash if it's a directory

		// Identify input/output files regardless of their original names
		var itemType string
		if strings.Contains(strings.ToLower(filename), "input") {
			itemType = "input.in" // Standardized name for input
		} else if strings.Contains(strings.ToLower(filename), "output") {
			itemType = "output.out" // Standardized name for output
		} else {
			// Skip files that are not clearly input or output
			log.Printf("Skipping file %s in directory '%s': not recognized as input or output.", filename, dir)
			continue
		}

		// Use the directory name as a temporary key to group files
		if _, ok := tempTestcasesByDir[dir]; !ok {
			tempTestcasesByDir[dir] = &struct {
				InputFile  *zip.File
				OutputFile *zip.File
			}{}
		}

		if itemType == "input.in" {
			// Check for duplicates within the same directory group
			if tempTestcasesByDir[dir].InputFile != nil {
				log.Printf("Warning: Duplicate input file found for directory '%s'. Keeping the first one encountered: %s", dir, tempTestcasesByDir[dir].InputFile.Name)
				continue // Skip this duplicate
			}
			tempTestcasesByDir[dir].InputFile = f
		} else if itemType == "output.out" {
			// Check for duplicates within the same directory group
			if tempTestcasesByDir[dir].OutputFile != nil {
				log.Printf("Warning: Duplicate output file found for directory '%s'. Keeping the first one encountered: %s", dir, tempTestcasesByDir[dir].OutputFile.Name)
				continue // Skip this duplicate
			}
			tempTestcasesByDir[dir].OutputFile = f
		}
	}

	// Now iterate through grouped files to create actual testcases, assigning sequential numbers
	var newTestcases []testCaseModel.TestCase
	currentTestcaseNumber := 1 // Start numbering from 1 for the final saved testcases

	for _, entry := range tempTestcasesByDir {
		if entry.InputFile == nil || entry.OutputFile == nil {
			// If a group (directory) doesn't have both input and output, log and skip it.
			// You might want more specific logging here to indicate which directory was incomplete.
			log.Printf("Skipping a testcase pair due to missing input or output file in its group.")
			continue
		}

		// Assign a new sequential testcase number
		testcaseNum := currentTestcaseNumber
		currentTestcaseNumber++

		// Create the target directory for the sequentially numbered testcase
		// This directory will be just the number (e.g., /testcases/<case_id>/1/)
		testcaseNumDir := filepath.Join(caseTestcasesDir, strconv.Itoa(testcaseNum))
		if err := os.MkdirAll(testcaseNumDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create testcase number directory %s: %w", testcaseNumDir, err)
		}

		// Handle input file: Save with standardized name
		inputFileName := fmt.Sprintf("t%03d.in", testcaseNum) // e.g., t001.in
		inputPath := filepath.Join(testcaseNumDir, inputFileName)
		inputURL := fmt.Sprintf("/private/testcases/%s/%d/%s", req.CaseID.String(), testcaseNum, inputFileName)

		if err := saveZipEntry(entry.InputFile, inputPath); err != nil {
			log.Printf("Failed to save input file for testcase %d (%s): %v", testcaseNum, entry.InputFile.Name, err)
			continue // Skip this testcase if input file fails to save
		}

		// Handle output file: Save with standardized name
		outputFileName := fmt.Sprintf("t%03d.out", testcaseNum) // e.g., t001.out
		outputPath := filepath.Join(testcaseNumDir, outputFileName)
		outputURL := fmt.Sprintf("/private/testcases/%s/%d/%s", req.CaseID.String(), testcaseNum, outputFileName)

		if err := saveZipEntry(entry.OutputFile, outputPath); err != nil {
			log.Printf("Failed to save output file for testcase %d (%s): %v", testcaseNum, entry.OutputFile.Name, err)
			continue // Skip this testcase if output file fails to save
		}

		newTestcases = append(newTestcases, testCaseModel.TestCase{
			CaseID:    req.CaseID,
			Number:    testcaseNum,
			InputUrl:  inputURL,
			OutputUrl: outputURL,
			CreatedAt: time.Now(),
		})
	}

	// Final validation: Ensure we processed at least one valid testcase pair
	if len(newTestcases) == 0 {
		return fmt.Errorf("no valid testcase input/output pairs found in the zip file. Ensure each pair is in a folder and files contain 'input' and 'output' in their names.")
	}

	// No need for the `tc.CreatedAt.IsZero()` check here as it's set when appended above.
	// The `InputURL` and `OutputURL` nil checks are also implicitly handled by `continue`
	// in the loop above if a pair is incomplete or fails to save.

	// 4. Save new testcases to DB in batch
	if err := s.testcaseRepo.SaveTestCaseBatch(ctx, newTestcases); err != nil {
		return fmt.Errorf("failed to save testcases to DB: %w", err)
	}

	return nil
}

// saveZipEntry is a helper function to open and save a zip.File to disk.
func saveZipEntry(f *zip.File, targetPath string) error {
	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetPath, err)
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open zip entry %s: %w", f.Name, err)
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return fmt.Errorf("failed to copy content for %s: %w", f.Name, err)
	}
	return nil
}

func (s testcaseServiceImpl) GetTestCasesByCaseID(ctx context.Context, caseID string) ([]responses.TestCaseResponse, error) {
	testcases, err := s.testcaseRepo.FindTestCaseByCaseID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to find test cases for case ID %s: %w", caseID, err)
	}

	resp := make([]responses.TestCaseResponse, len(testcases))
	for i, tc := range testcases {
		resp[i] = responses.TestCaseResponse{
			CaseID:    tc.CaseID.String(),
			Number:    tc.Number,
			InputUrl:  tc.InputUrl,
			OutputUrl: tc.OutputUrl,
		}
	}

	return resp, nil
}

func NewTestCaseService(testcaseRepo testCaseRepo.TestCaseRepository, caseRepo caseRepository.CaseRepository) TestCaseService {
	return &testcaseServiceImpl{
		testcaseRepo: testcaseRepo,
		caseRepo:     caseRepo,
	}
}
