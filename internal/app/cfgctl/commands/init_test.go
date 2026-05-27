package commands_test

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/FeelsCoderMan/cfgctl/internal/testutil"
)

func TestNewInitCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewInitCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"testKey": "testValue"}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("init_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "init", "--path", tmpPath, "--data", `{"testKey": "testValue"}`)

	if err != nil {
		t.Fatalf("init_test: expected no error, got %v", err)
	}

	if mockFileStorage.GetPath() != tmpPath {
		t.Fatalf("init_test: unexpected path: %s", mockFileStorage.GetPath())
	}

	if !mockFileStorage.IsDataValid(data) {
		t.Fatalf("init_test: data is not valid")
	}
}

func TestNewInitCmd_Success_EmptyData(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewInitCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("init_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "init", "--path", tmpPath, "--data", `{}`)

	if err != nil {
		t.Fatalf("init_test: expected no error, got %v", err)
	}

	if mockFileStorage.GetPath() != tmpPath {
		t.Fatalf("init_test: unexpected path: %s", mockFileStorage.GetPath())
	}

	if !mockFileStorage.IsDataValid(data) {
		t.Fatalf("init_test: data is not valid")
	}
}

func TestNewInitCmd_Success_MissingData(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewInitCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("init_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "init", "--path", tmpPath)

	if err != nil {
		t.Fatalf("init_test: expected no error, got %v", err)
	}

	if mockFileStorage.GetPath() != tmpPath {
		t.Fatalf("init_test: unexpected path: %s", mockFileStorage.GetPath())
	}

	if !mockFileStorage.IsDataValid(data) {
		t.Fatalf("init_test: data is not valid")
	}
}

func TestNewInitCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	c := commands.NewInitCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(c)

	err := testutil.RunCmd(rootCmd, "init", "--data", `{"testKey": "testValue"}`)

	if err == nil {
		t.Fatalf("init_test: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("init_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindMissingPath {
		t.Fatalf("init_test: expected MissingPath error, got %v", cmdError.ErrorKind)
	}
}

func TestNewInitCmd_Failure_InvalidData(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	c := commands.NewInitCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(c)
	tmpDir := testutil.CreatePackageTempDir(t)

	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"testKey": "testValue"}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("init_test: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "init", "--path", tmpPath, "--data", `{"testKey": "testValue":}`)

	if err == nil {
		t.Fatalf("init_test: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("init_test: expected CommandError, got %T (%v)", err, err)
	}

	if cmdError.ErrorKind != commands.KindInvalidJSON {
		t.Fatalf("init_test: expected InvalidJSON error, got %v", cmdError.ErrorKind)
	}

	var syntaxErr *json.SyntaxError
	if !errors.As(cmdError.Detail, &syntaxErr) {
		t.Fatalf("init_test: expected SyntaxError, got %T (%v)", cmdError.Detail, cmdError.Detail)
	}
}
