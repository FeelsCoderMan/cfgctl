package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func CreateTmpFile(path string, data any, shouldCorruptData bool) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dataByte, err := json.Marshal(data)
	if _, err := file.Write(dataByte); err != nil {
		return err
	}

	if shouldCorruptData {
		if _, err := file.Write([]byte("corrupt")); err != nil {
			return err
		}
	}

	return nil
}

func CreatePackageTempDir(t *testing.T) string {
	pkgDir, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	base := filepath.Join(pkgDir, "tmp")

	if err := os.MkdirAll(base, 0755); err != nil {
		t.Fatal(err)
	}

	dir, err := os.MkdirTemp(base, "test-")

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	return dir
}
