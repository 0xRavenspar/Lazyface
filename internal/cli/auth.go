package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// TokenData stores multiple tokens with an active one
type TokenData struct {
	ActiveToken string            `json:"active_token"`
	Tokens      map[string]string `json:"tokens"`
}

// getTokenFilePath returns the file path for storing tokens
func getTokenFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".lazyface", "tokens.json"), nil
}

// LoadTokens reads stored tokens or initializes a new struct if none exist
func LoadTokens() (*TokenData, error) {
	path, err := getTokenFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &TokenData{Tokens: make(map[string]string)}, nil
		}
		return nil, err
	}

	var tokens TokenData
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

// Save Tokens writes token data to disk
func SaveTokens(tokens *TokenData) error {
	path, err := getTokenFilePath()
	if err != nil {
		return err
	}

	//Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tokens, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func WhoAmI() (string, error) {
	return RunCommand("whoami")
}

func Login(token string, addGitCredential bool) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	args := []string{"login", "--token", token}
	if addGitCredential {
		args = append(args, "--add-to-git-credential")
	}

	output, err := RunCommand(args...)
	if err != nil {
		return err
	}

	var activeToken string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "The current active token is") {
			parts := strings.Split(line, ": ")
			if len(parts) > 1 {
				activeToken = strings.Trim(parts[1], "`")
			}
		}
	}

	if activeToken == "" {
		return fmt.Errorf("could not retrieve active token")
	}

	//Fetch username
	username, err := WhoAmI()
	if err != nil {
		return err
	}

	//Load and update token data
	tokens, err := LoadTokens()
	if err != nil {
		return err
	}

	tokens.Tokens[username] = token
	tokens.ActiveToken = activeToken

	//Save tokens
	if err := SaveTokens(tokens); err != nil {
		return err
	}

	fmt.Println("Login Successful")
	fmt.Println("Authentication details saved to ~/.lazyface/tokens.json")
	return nil
}

func Logout() error {
	_, err := RunCommand("logout")
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	tokenFilePath, err := getTokenFilePath()
	if err != nil {
		return fmt.Errorf("failed to get token file path: %w", err)
	}

	if err := os.Remove(tokenFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}
