package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_exists")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing file", tmpFile.Name(), true},
		{"nonexistent file", "/path/that/does/not/exist/file.txt", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.path); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpand(t *testing.T) {
	// Get home directory for testing
	home := os.Getenv("HOME")

	tests := []struct {
		name    string
		path    string
		wantErr bool
		check   func(string) bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: false,
			check:   func(s string) bool { return s == "" },
		},
		{
			name:    "path with tilde",
			path:    "~/test",
			wantErr: false,
			check:   func(s string) bool { return s == filepath.Join(home, "test") },
		},
		{
			name:    "absolute path",
			path:    "/tmp/test",
			wantErr: false,
			check:   func(s string) bool { return s == "/tmp/test" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.check(got) {
				t.Errorf("Expand() = %v, check failed", got)
			}
		})
	}
}

func TestChdir(t *testing.T) {
	// Get original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create temp directory
	tmpDir := t.TempDir()

	// Test successful directory change
	callbackExecuted := false
	err = Chdir(tmpDir, func() error {
		callbackExecuted = true
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		if wd != tmpDir {
			t.Errorf("Inside callback, wd = %v, want %v", wd, tmpDir)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Chdir() error = %v", err)
	}

	if !callbackExecuted {
		t.Error("Callback was not executed")
	}

	// Verify we're back to original directory
	currentWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if currentWd != originalWd {
		t.Errorf("After Chdir(), wd = %v, want %v", currentWd, originalWd)
	}
}

func TestChdirWithError(t *testing.T) {
	// Get original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()

	// Test with callback that returns error
	err = Chdir(tmpDir, func() error {
		return os.ErrNotExist
	})

	if err != os.ErrNotExist {
		t.Errorf("Chdir() error = %v, want %v", err, os.ErrNotExist)
	}

	// Verify we're still back to original directory even after error
	currentWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if currentWd != originalWd {
		t.Errorf("After Chdir() with error, wd = %v, want %v", currentWd, originalWd)
	}
}

func TestChdirInvalidPath(t *testing.T) {
	err := Chdir("/path/that/does/not/exist", func() error {
		return nil
	})

	if err == nil {
		t.Error("Chdir() with invalid path should return error")
	}
}

func TestGlob(t *testing.T) {
	// Create temp directory with files
	tmpDir := t.TempDir()

	// Create test files
	files := []string{"test1.txt", "test2.txt", "other.log"}
	for _, f := range files {
		file := filepath.Join(tmpDir, f)
		if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("match all txt files", func(t *testing.T) {
		var matched []string
		err := Glob(tmpDir, "*.txt", func(fileName string) error {
			matched = append(matched, filepath.Base(fileName))
			return nil
		})

		if err != nil {
			t.Errorf("Glob() error = %v", err)
		}

		if len(matched) != 2 {
			t.Errorf("Glob() matched %d files, want 2", len(matched))
		}
	})

	t.Run("callback returns error", func(t *testing.T) {
		count := 0
		err := Glob(tmpDir, "*.txt", func(fileName string) error {
			count++
			if count > 1 {
				return os.ErrPermission
			}
			return nil
		})

		if err != os.ErrPermission {
			t.Errorf("Glob() error = %v, want %v", err, os.ErrPermission)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		var matched []string
		err := Glob(tmpDir, "*.pdf", func(fileName string) error {
			matched = append(matched, fileName)
			return nil
		})

		if err != nil {
			t.Errorf("Glob() error = %v", err)
		}

		if len(matched) != 0 {
			t.Errorf("Glob() matched %d files, want 0", len(matched))
		}
	})
}
