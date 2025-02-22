package cli

import (
	"fmt"
)

// UploadToHuggingFace uploads files to a Hugging Face repository with advanced options.
// It uses a wrapper function for executing Hugging Face CLI commands.
func UploadToHuggingFace(repoID, localPath, pathInRepo string, repoType string, includePattern, excludePattern, deletePattern, commitMessage, revision string) (string, error) {
	// Validate that repoID, localPath, and pathInRepo are not empty
	if repoID == "" || localPath == "" || pathInRepo == "" {
		return "", fmt.Errorf("repo ID, local path, and path in repo cannot be empty")
	}

	// Build the command arguments
	args := []string{"upload", repoID, localPath, pathInRepo}

	// Add optional flags if provided
	if repoType != "" {
		args = append(args, "--repo-type", repoType)
	}
	if includePattern != "" {
		args = append(args, "--include", includePattern)
	}
	if excludePattern != "" {
		args = append(args, "--exclude", excludePattern)
	}
	if deletePattern != "" {
		args = append(args, "--delete", deletePattern)
	}
	if commitMessage != "" {
		args = append(args, "--commit-message", commitMessage)
	}
	if revision != "" {
		args = append(args, "--revision", revision)
	}

	// Use your wrapper function to execute the Hugging Face CLI command
	_, err := RunCommand(args...)
	if err != nil {
		return "", fmt.Errorf("failed to upload files to Hugging Face: %v", err)
	}

	return "SUCCESS", nil
}
