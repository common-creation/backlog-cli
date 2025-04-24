package main

import (
	"fmt"
	"os"

	"github.com/common-creation/backlog-cli/internal/client"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "backlog",
		Usage: "Backlog CLI tool",
		Commands: []*cli.Command{
			{
				Name:  "issue",
				Usage: "Manage Backlog issues",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List issues",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "project",
								Aliases: []string{"p"},
								Usage:   "Project key",
							},
							&cli.StringFlag{
								Name:    "status",
								Aliases: []string{"s"},
								Usage:   "Status ID",
							},
							&cli.IntFlag{
								Name:    "count",
								Aliases: []string{"c"},
								Value:   20,
								Usage:   "Number of issues to retrieve",
							},
						},
						Action: func(c *cli.Context) error {
							backlogClient, err := client.NewClient()
							if err != nil {
								return err
							}
							
							projectKey := c.String("project")
							statusID := c.String("status")
							count := c.Int("count")
							
							issues, err := backlogClient.ListIssues(projectKey, statusID, count)
							if err != nil {
								return err
							}
							
							for _, issue := range issues {
								fmt.Printf("#%d %s - %s\n", issue.ID, issue.IssueKey, issue.Summary)
							}
							
							return nil
						},
					},
					{
						Name:  "get",
						Usage: "Get issue details",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "key",
								Aliases:  []string{"k"},
								Usage:    "Issue key (e.g., PROJECT-123)",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							backlogClient, err := client.NewClient()
							if err != nil {
								return err
							}
							
							issueKey := c.String("key")
							
							issue, err := backlogClient.GetIssue(issueKey)
							if err != nil {
								return err
							}
							
							fmt.Printf("Issue: %s\n", issue.IssueKey)
							fmt.Printf("Summary: %s\n", issue.Summary)
							fmt.Printf("Status: %s\n", issue.Status.Name)
							fmt.Printf("Assignee: %s\n", issue.Assignee.Name)
							fmt.Printf("Description:\n%s\n", issue.Description)
							
							return nil
						},
					},
					{
						Name:  "create",
						Usage: "Create a new issue",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project",
								Aliases:  []string{"p"},
								Usage:    "Project key",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "summary",
								Aliases:  []string{"s"},
								Usage:    "Issue summary",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "description",
								Aliases: []string{"d"},
								Usage:   "Issue description",
							},
							&cli.IntFlag{
								Name:    "issue-type",
								Aliases: []string{"t"},
								Usage:   "Issue type ID",
								Required: true,
							},
							&cli.IntFlag{
								Name:    "priority",
								Aliases: []string{"pr"},
								Usage:   "Priority ID",
								Value:   3, // Normal priority
							},
						},
						Action: func(c *cli.Context) error {
							backlogClient, err := client.NewClient()
							if err != nil {
								return err
							}
							
							projectKey := c.String("project")
							summary := c.String("summary")
							description := c.String("description")
							issueTypeID := c.Int("issue-type")
							priorityID := c.Int("priority")
							
							issue, err := backlogClient.CreateIssue(projectKey, summary, description, issueTypeID, priorityID)
							if err != nil {
								return err
							}
							
							fmt.Printf("Created issue: %s\n", issue.IssueKey)
							
							return nil
						},
					},
				},
			},
			{
				Name:  "project",
				Usage: "Manage Backlog projects",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List projects",
						Action: func(c *cli.Context) error {
							backlogClient, err := client.NewClient()
							if err != nil {
								return err
							}
							
							projects, err := backlogClient.ListProjects()
							if err != nil {
								return err
							}
							
							for _, project := range projects {
								fmt.Printf("%s - %s\n", project.ProjectKey, project.Name)
							}
							
							return nil
						},
					},
				},
			},
			{
				Name:  "config",
				Usage: "Configure Backlog CLI",
				Subcommands: []*cli.Command{
					{
						Name:  "init",
						Usage: "Initialize configuration",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "space",
								Aliases:  []string{"s"},
								Usage:    "Backlog space URL (e.g., https://yourspace.backlog.com)",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "api-key",
								Aliases:  []string{"k"},
								Usage:    "API key",
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "read-only",
								Aliases: []string{"r"},
								Usage:   "Set to read-only mode (default: true)",
								Value:   true,
							},
						},
						Action: func(c *cli.Context) error {
							space := c.String("space")
							apiKey := c.String("api-key")
							readOnly := c.Bool("read-only")
							
							err := client.SaveConfig(space, apiKey, readOnly)
							if err != nil {
								return err
							}
							
							modeStr := "read-only"
							if !readOnly {
								modeStr = "read-write"
							}
							fmt.Printf("Configuration saved successfully (mode: %s)\n", modeStr)
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
