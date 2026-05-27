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

func TestNewGetCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get_test: failed to create tmp file: %v", err)
	}

	key := "name"
	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, key)

	if err != nil {
		t.Fatalf("get_test: expected no error on execution of get command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get(key)

	if err != nil {
		t.Fatalf("get_test: expected no error, got %v", err)
	}

	if value != "Alan" {
		t.Fatalf("get_test: expected value 'Alan', got %v", value)
	}

}

func TestNewGetCmd_Success_UnpresentKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get_test: failed to create tmp file: %v", err)
	}

	key := "birthday"
	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, key)

	if err != nil {
		t.Fatalf("get_test: expected no error on execution of get command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get(key)

	if err != nil {
		t.Fatalf("get_test: expected no error, got %v", err)
	}

	if value != nil {
		t.Fatalf("get_test: expected nil value, got %v", value)
	}
}

func TestNewGetCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "get", "name")

	if err == nil {
		t.Fatalf("get_test: expected error on missing path, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get_test: expected CommandError, got %v", err)
	}

	if cmdError.ErrorKind != commands.KindMissingPath {
		t.Fatalf("get_test: expected MissingPath error, got %v", cmdError.ErrorKind)
	}
}

func TestNewGetCmd_Failure_MissingKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath)

	if err == nil {
		t.Fatalf("get_test: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingArgs {
		t.Fatalf("get_test: expected MissingArgs error, got %v", cmdError.ErrorKind)
	}
}

func TestNewGetCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("get_test: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get_test: expected CommandError, got %v", err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("get_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var pathError *os.PathError
	if !errors.As(cmdError.Detail, &pathError) {
		t.Fatalf("get_test: expected PathError error, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}

func TestNewGetCmd_Failure_InvalidJsonFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("get_test: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("get_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var syntaxErr *json.SyntaxError
	if !errors.As(cmdError.Detail, &syntaxErr) {
		t.Fatalf("get_test: expected SyntaxError, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}
