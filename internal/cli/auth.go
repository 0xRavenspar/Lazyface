package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type TokenData struct {
	ActiveToken string            `json:"active_token"`
	Tokens      map[string]string `json:"tokens"`
}
type UserData struct {
	Name        string   `json:"name"`
	FullName    string   `json:"fullname"`
	Orgs        []string `json:"orgs"`
	TokenName   string   `json:"tokenname"`
	Permissions []string `json:"permissions"`
}

func getUserFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".lazyface", "user.json"), nil
}

func getTokenFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".lazyface", "tokens.json"), nil
}

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

func FetchAndStoreUserData(token string) error {
	url := "https://huggingface.co/api/whoami-v2"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch user data: %s", resp.Status)
	}

	var rawResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return err
	}

	// Extract auth and accessToken
	authData, ok := rawResponse["auth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse auth data")
	}
	accessToken, ok := authData["accessToken"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse accessToken")
	}

	// Extract displayName from accessToken
	displayName, ok := accessToken["displayName"].(string)
	if !ok {
		return fmt.Errorf("failed to parse displayName")
	}

	// Extract permissions directly from fineGrained.scoped
	var permissions []string
	if fineGrained, ok := accessToken["fineGrained"].(map[string]interface{}); ok {
		if scoped, ok := fineGrained["scoped"].([]interface{}); ok {
			for _, scope := range scoped {
				if scopeMap, ok := scope.(map[string]interface{}); ok {
					if perms, ok := scopeMap["permissions"].([]interface{}); ok {
						for _, perm := range perms {
							if permStr, ok := perm.(string); ok {
								permissions = append(permissions, permStr)
							}
						}
					}
				}
			}
		}
	}

	// Map extracted data
	userData := map[string]interface{}{
		"name":        rawResponse["name"],
		"fullname":    rawResponse["fullname"],
		"orgs":        rawResponse["orgs"],
		"tokenname":   displayName, // Save the token name
		"permissions": permissions,
	}

	// Save to file
	path, err := getUserFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadUserData() (*UserData, error) {
	path, err := getUserFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var userData UserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return nil, err
	}

	return &userData, nil
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

	err = FetchAndStoreUserData(token)
	if err != nil {
		return fmt.Errorf("login Successful but failed tp fetch user details: %w", err)
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

	lazyfaceDir, err := getUserFilePath()
	if err != nil {
		return fmt.Errorf("failed to get token file path: %w", err)
	}

	if err := os.RemoveAll(lazyfaceDir); err != nil {
		return fmt.Errorf("failed to delete user data: %w", err)
	}
	// tokenFilePath, err := getTokenFilePath()
	// if err != nil {
	// 	return fmt.Errorf("failed to get token file path: %w", err)
	// }
	//
	// if err := os.Remove(tokenFilePath); err != nil && !os.IsNotExist(err) {
	// 	return fmt.Errorf("failed to delete token file: %w", err)
	// }
	//
	return nil
}
