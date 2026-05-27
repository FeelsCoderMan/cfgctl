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

func TestNewListCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"testKey": "testValue"}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("list_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err != nil {
		t.Fatalf("list_test: expected no error, got %v", err)
	}

	if mockFileStorage.GetPath() != tmpPath {
		t.Fatalf("list_test: unexpected path: %s", mockFileStorage.GetPath())
	}

	if !mockFileStorage.IsDataValid(data) {
		t.Fatalf("list_test: data is not valid")
	}
}

func TestNewListCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	c := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(c)

	err := testutil.RunCmd(rootCmd, "list")

	if err == nil {
		t.Fatalf("list_test: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingPath {
		t.Fatalf("list_test: expected MissingPath error, got %v", cmdError.ErrorKind)
	}
}

func TestNewListCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err == nil {
		t.Fatalf("list_test: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("list_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var pathError *os.PathError
	if !errors.As(cmdError.Detail, &pathError) {
		t.Fatalf("list_test: expected PathError error, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}

func TestNewListCmd_Failure_InvalidFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("list_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err == nil {
		t.Fatalf("list_test: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindFileStoreLoad {
		t.Fatalf("list_test: expected FileStoreLoad error, got %v", cmdError.ErrorKind)
	}

	var syntaxerr *json.SyntaxError
	if !errors.As(cmdError.Detail, &syntaxerr) {
		t.Fatalf("list_test: expected JSON error, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}
