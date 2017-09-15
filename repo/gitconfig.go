package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"github.com/google/concourse-resources/internal/temputil"
)

var (
	execGit = realExecGit
)

type gitConfigManager struct {
	TempFileManager temputil.TempFileManager
	oldValues map[string]string
}

func (m *gitConfigManager) set(key, value string) error {
	if m.oldValues == nil {
		m.oldValues = make(map[string]string)
	}
	if _, exists := m.oldValues[key]; !exists {
		value, err := gitConfigGet(key)
		if err != nil {
			return err
		}
		m.oldValues[key] = value
	}
	return gitConfigSet(key, value)
}

func (m *gitConfigManager) setFileContents(key, fileContents string) error {
	// TODO: derive prefix from key
	tempFile, err := m.TempFileManager.Create("", fileContents)
	if err != nil {
		return err
	}
	return m.set(key, tempFile)
}

func (m *gitConfigManager) cleanup() {
	m.TempFileManager.Cleanup()
	for key, oldValue := range m.oldValues {
		err := gitConfigSet(key, oldValue)
		if err != nil {
			log.Printf("error resetting git config: %v", err)
		}
	}
	m.oldValues = nil
}

func realExecGit(args ...string) ([]byte, error) {
	return exec.Command("git", args...).Output()
}

func gitConfig(args ...string) (string, error, int) {
	args = append([]string{"config", "--global", "--null"}, args...)
	output, err := execGit(args...)

	exitStatus := 0
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			err = fmt.Errorf(
				"error running git %q\nstdout:\n%s\nstderr:\n%s",
				args, output, exitErr.Stderr)

			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitStatus = status.ExitStatus()
			}
		}
	}
	return strings.TrimSuffix(string(output), "\x00"), err, exitStatus
}

func gitConfigGet(key string) (string, error) {
		output, err, exitStatus := gitConfig(key)
		if err != nil {
			if exitStatus == 1 {
				err = nil // the config key isn't set; not an error
			}
			return "", err
		}
		return output, nil
}

func gitConfigSet(key, value string) error {
	_, err, _ := gitConfig(key, value)
	return err
}
