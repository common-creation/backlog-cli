package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kenzo0107/backlog"
)

type Config struct {
	Space     string `json:"space"`
	APIKey    string `json:"api_key"`
	ReadOnly  bool   `json:"read_only"`
}

type Client struct {
	backlogClient *backlog.Client
	readOnly      bool
}

func NewClient() (*Client, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	c := backlog.New(config.Space, config.APIKey)
	return &Client{
		backlogClient: c,
		readOnly:      config.ReadOnly,
	}, nil
}

func SaveConfig(space, apiKey string, readOnly bool) error {
	config := Config{
		Space:     space,
		APIKey:    apiKey,
		ReadOnly:  readOnly,
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func loadConfig() (*Config, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(configDir, "config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".backlog-cli"), nil
}

func (c *Client) ListIssues(projectKey, statusID string, count int) ([]*backlog.Issue, error) {
	countPtr := backlog.Int(count)
	options := &backlog.GetIssuesOptions{
		Count: countPtr,
	}

	if projectKey != "" {
		project, err := c.backlogClient.GetProject(projectKey)
		if err == nil && project != nil && project.ID != nil {
			options.ProjectIDs = []int{*project.ID}
		}
	}

	if statusID != "" {
		var statusIDInt int
		fmt.Sscanf(statusID, "%d", &statusIDInt)
		if statusIDInt > 0 {
			options.StatusIDs = []int{statusIDInt}
		}
	}

	issues, err := c.backlogClient.GetIssues(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	return issues, nil
}

func (c *Client) GetIssue(issueKey string) (*backlog.Issue, error) {
	issue, err := c.backlogClient.GetIssue(issueKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return issue, nil
}

func (c *Client) CreateIssue(projectKey, summary, description string, issueTypeID, priorityID int) (*backlog.Issue, error) {
	if c.readOnly {
		return nil, fmt.Errorf("cannot create issue: client is in read-only mode")
	}

	var projectID int
	project, err := c.backlogClient.GetProject(projectKey)
	if err == nil && project != nil && project.ID != nil {
		projectID = *project.ID
	} else {
		fmt.Sscanf(projectKey, "%d", &projectID)
	}

	input := &backlog.CreateIssueInput{
		ProjectID:   backlog.Int(projectID),
		Summary:     backlog.String(summary),
		IssueTypeID: backlog.Int(issueTypeID),
		PriorityID:  backlog.Int(priorityID),
	}

	if description != "" {
		input.Description = backlog.String(description)
	}

	issue, err := c.backlogClient.CreateIssue(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return issue, nil
}

func (c *Client) ListProjects() ([]*backlog.Project, error) {
	projects, err := c.backlogClient.GetProjects(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, nil
}

func (c *Client) GetIssueComments(issueKey string) ([]*backlog.IssueComment, error) {
	comments, err := c.backlogClient.GetIssueComments(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue comments: %w", err)
	}

	return comments, nil
}
