package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/lib"
	"github.com/exceptionate/gitaur/internal/models"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
)

func showExperiences(short bool) {
	rows, err := db.Conn.Query(`
		SELECT id, title, company, startDate, endDate, description, tags
		FROM experiences
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var experiences models.Experiences

	for rows.Next() {
		var experience models.Experience
		var tagsJSON string

		err := rows.Scan(
			&experience.ID,
			&experience.Title,
			&experience.Company,
			&experience.StartDate,
			&experience.EndDate,
			&experience.Description,
			&tagsJSON,
		)
		if err != nil {
			panic(err)
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &experience.Tags); err != nil {
				experience.Tags = nil
			}
		}

		experiences = append(experiences, experience)
	}

	if len(experiences) == 0 {
		fmt.Println(ui.Info.Render("No experiences found."))
		return
	}

	fmt.Printf("%s\n\n", ui.Info.Render(fmt.Sprintf("%d experiences found:", len(experiences))))

	if short {
		for _, experience := range experiences {
			fmt.Printf("%s: %s\n", ui.Label.Render(strconv.Itoa(experience.ID)), experience.Title)
		}
	} else {
		for i, experience := range experiences {
			experienceInfo := fmt.Sprintf(
				"%s %d\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
				ui.Label.Render("ID:"), experience.ID,
				ui.Label.Render("Title:"), experience.Title,
				ui.Label.Render("Company:"), experience.Company,
				ui.Label.Render("Start Date:"), experience.StartDate,
				ui.Label.Render("End Date:"), experience.EndDate,
				ui.Label.Render("Description:"), experience.Description,
				ui.Label.Render("Tags:"), strings.Join(experience.Tags, ", "),
			)

			fmt.Println(experienceInfo)

			if i != len(experiences)-1 {
				fmt.Println("----------------------------")
			}
		}
	}
}

func addExperience() {
	fmt.Println(ui.Info.Render("Add Experience"))

	var experience models.Experience

	experience.Title = lib.PromptWithDefault("Title", "")
	experience.Company = lib.PromptWithDefault("Company", "")
	experience.StartDate = lib.PromptWithDefault("Start date", "")
	experience.EndDate = lib.PromptWithDefault("End date", "")
	experience.Description = lib.PromptWithDefault("Description", "")
	tagsInput := lib.PromptWithDefault("Tags (comma separated)", "")
	if tagsInput != "" {
		experience.Tags = strings.Split(tagsInput, ",")
	}

	tagsJSON, err := json.Marshal(experience.Tags)
	if err != nil {
		panic(err)
	}

	_, err = db.Conn.Exec(`
		INSERT INTO experiences (
			title,
			company,
			startDate,
			endDate,
			description,
			tags
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		experience.Title,
		experience.Company,
		experience.StartDate,
		experience.EndDate,
		experience.Description,
		string(tagsJSON),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(ui.Success.Render("Experience added successfully."))
}

func updateExperience(id int, field string) {
	if strings.EqualFold(field, "id") || !lib.HasField[models.Experience](field) {
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
			"SELECT %s FROM experiences WHERE id = ? LIMIT 1",
			field,
		),
		id,
	).Scan(&currentValue)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(
				ui.Error.Render(
					"Experience not found.",
				),
			)
			return
		}
		panic(err)
	}

	value := ""
	if strings.EqualFold(field, "tags") {
		var tags []string
		if currentValue.String != "" {
			_ = json.Unmarshal([]byte(currentValue.String), &tags)
		}
		value = lib.PromptWithDefault(
			"Update tags (comma separated)",
			strings.Join(tags, ","),
		)
		if value != "" {
			tags = strings.Split(value, ",")
		}
		encoded, err := json.Marshal(tags)
		if err != nil {
			panic(err)
		}
		value = string(encoded)
	} else {
		value = lib.PromptWithDefault(
			fmt.Sprintf(
				"Update %s",
				field,
			),
			currentValue.String,
		)
	}

	_, err = db.Conn.Exec(
		fmt.Sprintf(
			"UPDATE experiences SET %s = ? WHERE id = ?",
			field,
		),
		value,
		id,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		ui.Success.Render(
			fmt.Sprintf(
				"Experience %d's %s updated successfully.",
				id,
				field,
			),
		),
	)
}

func deleteExperience(id int) {
	result, err := db.Conn.Exec(
		"DELETE FROM experiences WHERE id = ?",
		id,
	)

	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowsAffected == 0 {
		fmt.Println(
			ui.Error.Render(
				"Experience not found.",
			),
		)
		return
	}

	fmt.Println(
		ui.Success.Render(
			fmt.Sprintf(
				"Experience '%d' deleted successfully.",
				id,
			),
		),
	)
}

func experiences(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showExperiences(false)
		return
	}

	if len(args) == 1 && args[0] == "--id" {
		showExperiences(true)
		return
	}

	if len(args) == 1 && args[0] == "add" {
		addExperience()
		return
	}

	if len(args) == 3 && (args[1] == "--set" || args[1] == "-s") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid experience id. Use a numeric id."))
			return
		}
		updateExperience(
			id,
			args[2],
		)
		return
	}

	if len(args) == 3 {
		if (args[1] == "--delete" || args[1] == "-d") && (args[2] == "--force" || args[2] == "-f") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid experience id. Use a numeric id."))
				return
			}
			deleteExperience(id)
			return
		}

		if (args[1] == "--force" || args[1] == "-f") && (args[2] == "--delete" || args[2] == "-d") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid experience id. Use a numeric id."))
				return
			}
			deleteExperience(id)
			return
		}
	}

	if len(args) == 2 && (args[1] == "--delete" || args[1] == "-d") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid experience id. Use a numeric id."))
			return
		}

		fmt.Printf(
			"%s Are you sure you want to delete the experience '%d'? This action cannot be undone. (y/N): ",
			ui.Warning.Render("Warning:"),
			id,
		)

		var confirmation string
		fmt.Scanln(&confirmation)

		if strings.ToLower(confirmation) != "y" {
			fmt.Println(ui.Info.Render("Experience deletion cancelled."))
			return
		}
		deleteExperience(
			id,
		)
		return
	}

	fmt.Println(ui.Error.Render("Invalid command."))
}

var experiencesCmd = &cobra.Command{
	Use:                "experiences",
	Aliases:            []string{"exp", "exps"},
	Short:              "List and manage your experiences",
	Run:                experiences,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(experiencesCmd)
}
