package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://huggingface.co/api"
)

// CreateRepo creates a repository on Hugging Face.
func CreateRepo(hfToken, repoType, repoName, organization string, isPrivate bool, sdk string) error {
	// Validate repoType
	if repoType != "model" && repoType != "dataset" && repoType != "space" {
		return errors.New("invalid repo type. Must be 'model', 'dataset', or 'space'")
	}

	// Prepare payload
	payload := map[string]interface{}{
		"type":    repoType,
		"name":    repoName,
		"private": isPrivate,
		"sdk":     sdk,
	}
	if organization != "" {
		payload["organization"] = organization
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repos/create", baseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hfToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create repo: %s", string(body))
	}

	return nil
}

// DeleteRepo deletes a repository on Hugging Face.
func DeleteRepo(hfToken, repoType, repoName, organization string) error {
	// Validate repoType
	if repoType != "model" && repoType != "dataset" && repoType != "space" {
		return errors.New("invalid repo type. Must be 'model', 'dataset', or 'space'")
	}

	// Prepare payload
	payload := map[string]string{
		"type": repoType,
		"name": repoName,
	}
	if organization != "" {
		payload["organization"] = organization
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/repos/delete", baseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hfToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete repo: %s", string(body))
	}

	return nil
}

// UpdateRepoVisibility updates the visibility of a repository.
func UpdateRepoVisibility(hfToken, repoType, repoID string, isPrivate bool) error {
	// Validate repoType
	if repoType != "model" && repoType != "dataset" && repoType != "space" {
		return errors.New("invalid repo type. Must be 'model', 'dataset', or 'space'")
	}

	// Prepare payload
	payload := map[string]bool{
		"private": isPrivate,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/repos/%s/%s/settings", baseURL, repoType, repoID), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hfToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update repo visibility: %s", string(body))
	}

	return nil
}

// MoveRepo moves or renames a repository.
func MoveRepo(hfToken, fromRepo, toRepo, repoType string) error {
	// Validate repoType
	if repoType != "model" && repoType != "dataset" && repoType != "space" {
		return errors.New("invalid repo type. Must be 'model', 'dataset', or 'space'")
	}

	// Prepare payload
	payload := map[string]string{
		"fromRepo": fromRepo,
		"toRepo":   toRepo,
		"type":     repoType,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repos/move", baseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hfToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to move repo: %s", string(body))
	}

	return nil
}
