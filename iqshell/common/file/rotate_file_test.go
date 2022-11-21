//go:build unit

package file

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func testDir() (string, error) {
	home, err := utils.GetHomePath()
	if err != nil {
		return home, err
	}
	return filepath.Join(home, "qshell_test", "tmp", "rotate_file"), nil
}

func TestNew(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Fatal("create tmp dir error", err)
	}
	_ = os.RemoveAll(dir)

	r, nErr := NewRotateFile(filepath.Join(dir, "test.txt"), RotateOptionMaxSize(10), RotateOptionMaxLine(10))
	if nErr != nil {
		t.Fatal(nErr)
	}

	if e := r.Close(); e != nil {
		t.Fatal(e)
	}

	_ = os.RemoveAll(dir)
}

func TestLineCount01(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Fatal("create tmp dir error", err)
	}
	_ = os.RemoveAll(dir)

	var maxLine int64 = 5
	r, nErr := NewRotateFile(filepath.Join(dir, "test.txt"), RotateOptionMaxLine(maxLine))
	if nErr != nil {
		t.Fatal(nErr)
	}

	for i := 0; i < 2; i++ {
		_, _ = r.Write([]byte(fmt.Sprintf("line 1: %d", i)))
		_, _ = r.Write([]byte(fmt.Sprintf("-%d\nline 2: %d\n", i, i)))
		_, _ = r.Write([]byte(fmt.Sprintf("line 3: %d\nline 4: %d\nline 5: %d\n", i, i, i)))
	}

	if e := r.Close(); e != nil {
		t.Fatal(e)
	}

	wErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		lineCount, lErr := utils.FileLineCounts(path)
		if lErr != nil {
			return lErr
		}
		if lineCount > maxLine {
			return errors.New("line count bigger than max line")
		}
		return nil
	})

	if wErr != nil {
		t.Fatal(wErr)
	}

	_ = os.RemoveAll(dir)
}

func TestLinesCount02(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Fatal("create tmp dir error", err)
	}
	_ = os.RemoveAll(dir)

	var maxLine int64 = 1
	r, nErr := NewRotateFile(filepath.Join(dir, "test.txt"), RotateOptionMaxLine(maxLine))
	if nErr != nil {
		t.Fatal(nErr)
	}

	for i := 0; i < 2; i++ {
		_, _ = r.Write([]byte(fmt.Sprintf("line 1: %d", i)))
		_, _ = r.Write([]byte(fmt.Sprintf("-%d\nline 2: %d\n", i, i)))
		_, _ = r.Write([]byte(fmt.Sprintf("line 3: %d\nline 4: %d", i, i)))
		_, _ = r.Write([]byte(fmt.Sprintf("\nline 5: %d\n", i)))
	}

	if e := r.Close(); e != nil {
		t.Fatal(e)
	}

	wErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		lineCount, lErr := utils.FileLineCounts(path)
		if lErr != nil {
			return lErr
		}
		if lineCount > maxLine {
			return errors.New("line count bigger than max line")
		}
		return nil
	})

	if wErr != nil {
		t.Fatal(wErr)
	}

	_ = os.RemoveAll(dir)
}

func TestLinesFileSize(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Fatal("create tmp dir error", err)
	}
	_ = os.RemoveAll(dir)

	var maxLine int64 = 2
	r, nErr := NewRotateFile(filepath.Join(dir, "test.txt"),
		RotateOptionMaxLine(maxLine),
		RotateOptionMaxSize(44),
		RotateOptionFileHeader("Key\tHash\tSize\tStatus\tModTime\tUser"))
	if nErr != nil {
		t.Fatal(nErr)
	}

	for i := 0; i < 2; i++ {
		_, _ = r.Write([]byte(fmt.Sprintf("line 1: %d", i)))
		_, _ = r.Write([]byte(fmt.Sprintf("-%d\nline 2: %d\n", i, i)))
		_, _ = r.Write([]byte(fmt.Sprintf("line 3: %d\nline 4: %d", i, i)))
		_, _ = r.Write([]byte(fmt.Sprintf("\nline 5: %d\n", i)))
	}

	if e := r.Close(); e != nil {
		t.Fatal(e)
	}

	wErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		lineCount, lErr := utils.FileLineCounts(path)
		if lErr != nil {
			return lErr
		}
		if lineCount > maxLine {
			return errors.New("line count bigger than max line")
		}
		return nil
	})

	if wErr != nil {
		t.Fatal(wErr)
	}

	_ = os.RemoveAll(dir)
}
