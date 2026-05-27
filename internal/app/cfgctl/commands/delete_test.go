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

func TestNewDeleteCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "name")

	if err != nil {
		t.Fatalf("delete_test: expected no error on execution of delete command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get("name")

	if value != nil {
		t.Fatalf("delete_test: expected nil value, got %v", value)
	}

	if err != nil {
		t.Fatalf("delete_test: expected no error, got %v", err)
	}
}

func TestNewDeleteCmd_Success_UnpresentKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "birthday")

	if err != nil {
		t.Fatalf("delete_test: expected no error on execution of delete command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get("birthday")

	if value != nil {
		t.Fatalf("delete_test: expected nil value, got %v", value)
	}

	if err != nil {
		t.Fatalf("delete_test: expected no error, got %v", err)
	}
}

func TestNewDeleteCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "name")

	if err == nil {
		t.Fatalf("delete_test: expected error on missing path, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingPath {
		t.Fatalf("delete_test: expected MissingPath error, got %v", cmdError.ErrorKind)
	}
}

func TestNewDeleteCmd_Failure_MissingKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath)

	if err == nil {
		t.Fatalf("delete_test: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingArgs {
		t.Fatalf("delete_test: expected MissingArgs error, got %v", cmdError.ErrorKind)
	}
}

func TestNewDeleteCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("delete_test: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("delete_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var pathError *os.PathError
	if !errors.As(cmdError.Detail, &pathError) {
		t.Fatalf("delete_test: expected PathError error, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}

func TestNewDeleteCmd_Failure_InvalidJsonFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("delete_test: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("delete_test: expected file store load error, got %v", cmdError.ErrorKind)
	}

	var syntaxErr *json.SyntaxError
	if !errors.As(cmdError.Detail, &syntaxErr) {
		t.Fatalf("delete_test: expected SyntaxError, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}
