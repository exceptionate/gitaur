package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/lib"
	"github.com/exceptionate/gitaur/internal/models"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
)

func showProjects() {
	rows, err := db.Conn.Query(`
		SELECT name, repo, owner, url, tech, desc, text, startDate, endDate
		FROM projects
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var projects models.Projects

	for rows.Next() {
		var project models.Project
		var tech string

		err := rows.Scan(
			&project.Name,
			&project.Repo,
			&project.Owner,
			&project.Url,
			&tech,
			&project.Desc,
			&project.Text,
			&project.StartDate,
			&project.EndDate,
		)
		if err != nil {
			panic(err)
		}

		if tech != "" {
			project.Tech = strings.Split(tech, ",")
		}

		projects = append(projects, project)
	}

	if len(projects) == 0 {
		fmt.Println(ui.Info.Render("No projects found."))
		return
	}

	fmt.Printf("%s\n\n", ui.Info.Render(fmt.Sprintf("%d projects found:", len(projects))))

	for i, project := range projects {
		projectInfo := fmt.Sprintf(
			"%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
			ui.Label.Render("Name:"), project.Name,
			ui.Label.Render("Repo:"), project.Repo,
			ui.Label.Render("Owner:"), project.Owner,
			ui.Label.Render("URL:"), project.Url,
			ui.Label.Render("Tech:"), strings.Join(project.Tech, ", "),
			ui.Label.Render("Description:"), project.Desc,
			ui.Label.Render("Start Date:"), project.StartDate,
			ui.Label.Render("End Date:"), project.EndDate,
		)

		fmt.Println(projectInfo)

		if i != len(projects)-1 {
			fmt.Println("----------------------------")
		}
	}
}

func addProject() {
	var user models.User

	err := db.Conn.QueryRow(
		"SELECT username FROM user LIMIT 1",
	).Scan(&user.Username)
	if err != nil {
		panic(err)
	}

	fmt.Println(ui.Info.Render("Add Project"))

	repo := lib.PromptWithDefault(
		fmt.Sprintf("%s (https://github.com/owner_name/%s)", ui.Label.Render("Repository Name"), ui.Info.Render("repository_name")),
		"",
	)

	owner := lib.PromptWithDefault(
		fmt.Sprintf("%s (https://github.com/%s/repository_name)", ui.Label.Render("Owner Name"), ui.Info.Render("owner_name")),
		user.Username,
	)

	res, err := lib.GithubRequest(
		fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var repository struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
		HTMLURL     string `json:"html_url"`
		HasPages    bool   `json:"has_pages"`
	}

	json.NewDecoder(res.Body).Decode(&repository)

	var project models.Project

	project.Name = lib.PromptWithDefault(
		"Project name",
		repository.Name,
	)
	project.Repo = repo
	project.Owner = owner

	defaultUrl := repository.HTMLURL

	if repository.HasPages && repository.Homepage != "" {
		defaultUrl = repository.Homepage
	}

	project.Url = lib.PromptWithDefault(
		"Project URL",
		defaultUrl,
	)

	res, err = lib.GithubRequest(
		fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", owner, repo),
	)

	if err == nil {
		defer res.Body.Close()

		var languages map[string]int

		json.NewDecoder(res.Body).Decode(&languages)

		for language := range languages {
			project.Tech = append(project.Tech, language)
		}
	}

	techInput := lib.PromptWithDefault(
		"Technologies (comma separated)",
		strings.Join(project.Tech, ", "),
	)

	project.Tech = strings.Split(techInput, ",")

	res, err = lib.GithubRequest(
		fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", owner, repo),
	)

	readmeContent := ""

	if err == nil {
		defer res.Body.Close()

		var readme struct {
			Content string `json:"content"`
		}

		json.NewDecoder(res.Body).Decode(&readme)

		readmeContent = lib.DecodeGithubContent(
			readme.Content,
		)
	}

	project.Text = repository.Description + "\n\n" + readmeContent

	project.Desc = lib.PromptWithDefault(
		"Short description",
		repository.Description,
	)

	res, err = lib.GithubRequest(
		fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo),
	)

	if err == nil {
		defer res.Body.Close()

		var commits []struct {
			Commit struct {
				Author struct {
					Date string `json:"date"`
				} `json:"author"`
			} `json:"commit"`
		}

		json.NewDecoder(res.Body).Decode(&commits)

		if len(commits) > 0 {
			start, _ := time.Parse(
				time.RFC3339,
				commits[len(commits)-1].Commit.Author.Date,
			)

			project.StartDate = start.Format("2006-01-02")

			last, _ := time.Parse(
				time.RFC3339,
				commits[0].Commit.Author.Date,
			)

			if time.Since(last).Hours() > 24*7 {
				project.EndDate = last.Format("2006-01-02")
			}
		}
	}

	project.StartDate = lib.PromptWithDefault(
		"Start date",
		project.StartDate,
	)

	project.EndDate = lib.PromptWithDefault(
		"End date",
		project.EndDate,
	)

	_, err = db.Conn.Exec(`
		INSERT INTO projects (
			name,
			repo,
			owner,
			url,
			tech,
			desc,
			text,
			startDate,
			endDate
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		project.Name,
		project.Repo,
		project.Owner,
		project.Url,
		strings.Join(project.Tech, ","),
		project.Desc,
		project.Text,
		project.StartDate,
		project.EndDate,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(ui.Success.Render("Project added successfully."))
}

func updateProject(repo string, field string) {
	if !lib.HasField[models.Project](field) {
		fmt.Println(
			ui.Error.Render(
				"Invalid field name.",
			),
		)
		return
	}

	var currentValue sql.NullString

	err := db.Conn.QueryRow(
		fmt.Sprintf(
			"SELECT %s FROM projects WHERE repo = ? LIMIT 1",
			field,
		),
		repo,
	).Scan(&currentValue)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(
				ui.Error.Render(
					"Project not found.",
				),
			)
			return
		}
		panic(err)

	}

	value := lib.PromptWithDefault(
		fmt.Sprintf(
			"Update %s",
			field,
		),
		currentValue.String,
	)

	_, err = db.Conn.Exec(
		fmt.Sprintf(
			"UPDATE projects SET %s = ? WHERE repo = ?",
			field,
		),
		value,
		repo,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		ui.Success.Render(
			fmt.Sprintf(
				"%s updated successfully.",
				field,
			),
		),
	)
}

func projects(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showProjects()
		return
	}

	if len(args) == 1 && args[0] == "add" {
		addProject()
		return
	}

	if len(args) == 3 && args[1] == "set" {
		updateProject(
			args[0],
			args[2],
		)
		return
	}

	fmt.Println(ui.Error.Render("Invalid command. Use 'gitaur projects' to list projects, 'gitaur projects add' to add a new project, or 'gitaur projects <repo> set <field>' to update a project field."))
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List and manage your projects",
	Run:   projects,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
