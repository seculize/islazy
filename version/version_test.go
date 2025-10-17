package version

import (
	"regexp"
	"testing"
)

func TestVersionFormat(t *testing.T) {
	// Check that version is not empty
	if Version == "" {
		t.Error("Version should not be empty")
	}

	// Check that version follows semantic versioning pattern (e.g., 1.10.6)
	semverPattern := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !semverPattern.MatchString(Version) {
		t.Errorf("Version %q does not follow semantic versioning format", Version)
	}
}

func TestVersionConstant(t *testing.T) {
	// Just verify the version constant is accessible
	v := Version
	if v == "" {
		t.Error("Version constant is empty")
	}
}
