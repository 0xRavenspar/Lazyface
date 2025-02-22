package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type RepoInfo struct {
	Siblings []struct { //siblings is an array of objects
		Filename string `json:"rfilename"`
	} `json:"siblings"`
}

func ListRepoFiles(repoID string) ([]string, error) {
	url := fmt.Sprintf("https://huggingface.co/api/models/%s", repoID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var repoInfo RepoInfo
	if err := json.Unmarshal(body, &repoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	//Extract filename
	files := make([]string, len(repoInfo.Siblings))
	for i, file := range repoInfo.Siblings {
		files[i] = file.Filename
	}

	return files, nil
}

func GetDownloadPath(choice string, userPath string) (string, error) {
	switch choice {
	case "default":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		downloadPath := filepath.Join(homeDir, "Downloads", "hfmodels")
		return downloadPath, nil

	case "pwd":
		return os.Getwd()
	case "custom":
		if userPath == "" {
			return "", fmt.Errorf("no custom path provided")
		}
		return createDirectory(userPath)
	default:
		return "", fmt.Errorf("invalid choice")
	}
}

func createDirectory(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to createDirectory: %w", err)
		}
	}
	return path, nil
}

func Download(repoID string, files []string, downloadPath string, updateStatus func(string), updateProgress func(float64)) error {
	noFiles := len(files)
	for i, file := range files {
		updateStatus(fmt.Sprintf("Downloading File %d/%d: %s", i+1, len(files), file))

		args := []string{"download", repoID, "--local-dir", downloadPath, "--include", file}
		cmd := exec.Command("huggingface-cli", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to download file %s: %v\nOutput: %s", file, err, string(output))
		}
		progress := (float64(i+1) / float64(noFiles))
		updateProgress(progress)
	}

	updateStatus(fmt.Sprintf("Download complete! %d files downloaded.", len(files)))
	updateProgress(1.0)
	return nil
}

