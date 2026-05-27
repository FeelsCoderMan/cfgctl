package commands_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/FeelsCoderMan/cfgctl/internal/testutil"
)

func TestNewSetCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set_test: failed to create tmp file: %v", err)
	}

	key := "name"
	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, key, "Bob")

	if err != nil {
		t.Fatalf("set_test: expected no error on execution of set command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get(key)

	if err != nil {
		t.Fatalf("set_test: expected no error, got %v", err)
	}

	if value != "Bob" {
		t.Fatalf("set_test: expected value 'Bob', got %v", value)
	}
}

func TestNewSetCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "set", "name", "Bob")

	if err == nil {
		t.Fatalf("set_test: expected error on missing path, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingPath {
		t.Fatalf("set_test: expected MissingPath error, got %v", cmdError.ErrorKind)
	}
}

func TestNewSetCmd_Failure_MissingKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, "Bob")

	if err == nil {
		t.Fatalf("set_test: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingArgs {
		t.Fatalf("set_test: expected MissingArgs error, got %v", cmdError.ErrorKind)
	}
}

func TestNewSetCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, "name", "Bob")

	if err == nil {
		t.Fatalf("set_test: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("set_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var pathError *os.PathError
	if !errors.As(cmdError.Detail, &pathError) {
		t.Fatalf("set_test: expected PathError error, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}

func TestNewSetCmd_Failure_InvalidFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, "name", "Bob")

	if err == nil {
		t.Fatalf("set_test: expected error on invalid file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("set_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var syntaxErr *json.SyntaxError
	if !errors.As(cmdError.Detail, &syntaxErr) {
		t.Fatalf("set_test: expected SyntaxError, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}
