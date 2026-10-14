package security

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestEmsCoreDoesNotImportEducationDomains enforces the EMS Core dependency
// rule: no file under backend/domain/ems may import an Education domain package.
// The Core must remain usable even if the Education addon does not exist.
func TestEmsCoreDoesNotImportEducationDomains(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve the test file path")
	}

	backendRoot := filepath.Dir(filepath.Dir(thisFile))
	emsRoot := filepath.Join(backendRoot, "domain", "ems")

	forbidden := []string{
		"domain/school",
		"domain/campus",
		"domain/class",
		"domain/staff",
		"domain/student",
		"domain/admission",
		"domain/enrollment",
		"domain/invoice",
		"domain/payments",
		"domain/curriculum",
		"domain/assessment",
		"domain/timetable",
		"domain/lessonPlan",
		"domain/schemeOfWork",
		"domain/catering",
		"domain/reportSettings",
		"domain/reportTemplate",
		"domain/promotion",
		"domain/interview",
		"domain/attendance",
		"domain/parent",
		"domain/item",
		"domain/sourceOfSchoolInfo",
		"domain/previousSchool",
		"domain/curriculumdelivery",
		"domain/preAdmissionInterviewTest",
		"domain/progressReport",
		"domain/activityreporting",
		"domain/statistics",
		"domain/studentDetails",
		"domain/subject",
		"domain/academicYear",
		"domain/studyPeriod",
	}

	err := filepath.WalkDir(emsRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		for _, bad := range forbidden {
			if strings.Contains(string(content), bad) {
				t.Errorf("%s imports Education domain %q; EMS Core must not depend on Education", path, bad)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("walking %s: %v", emsRoot, err)
	}
}
