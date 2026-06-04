package cmd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/lib"
	"github.com/exceptionate/gitaur/internal/models"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
)

func showEducation(short bool) {
	rows, err := db.Conn.Query(`
		SELECT id, school, degree, field, startDate, endDate, grade, description
		FROM education
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var education models.EducationList

	for rows.Next() {
		var entry models.Education

		err := rows.Scan(
			&entry.ID,
			&entry.School,
			&entry.Degree,
			&entry.Field,
			&entry.StartDate,
			&entry.EndDate,
			&entry.Grade,
			&entry.Description,
		)
		if err != nil {
			panic(err)
		}

		education = append(education, entry)
	}

	if len(education) == 0 {
		fmt.Println(ui.Info.Render("No education entries found."))
		return
	}

	fmt.Printf("%s\n\n", ui.Info.Render(fmt.Sprintf("%d education entries found:", len(education))))

	if short {
		for _, entry := range education {
			fmt.Printf("%s: %s\n", ui.Label.Render(strconv.Itoa(entry.ID)), entry.School)
		}
	} else {
		for i, entry := range education {
			entryInfo := fmt.Sprintf(
				"%s %d\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
				ui.Label.Render("ID:"), entry.ID,
				ui.Label.Render("School:"), entry.School,
				ui.Label.Render("Degree:"), entry.Degree,
				ui.Label.Render("Field:"), entry.Field,
				ui.Label.Render("Start Date:"), entry.StartDate,
				ui.Label.Render("End Date:"), entry.EndDate,
				ui.Label.Render("Grade:"), entry.Grade,
				ui.Label.Render("Description:"), entry.Description,
			)

			fmt.Println(entryInfo)

			if i != len(education)-1 {
				fmt.Println("----------------------------")
			}
		}
	}
}

func addEducation() {
	fmt.Println(ui.Info.Render("Add Education"))

	var entry models.Education

	entry.School = lib.PromptWithDefault("School", "")
	entry.Degree = lib.PromptWithDefault("Degree", "")
	entry.Field = lib.PromptWithDefault("Field", "")
	entry.StartDate = lib.PromptWithDefault("Start date", "")
	entry.EndDate = lib.PromptWithDefault("End date", "")
	entry.Grade = lib.PromptWithDefault("Grade", "")
	entry.Description = lib.PromptWithDefault("Description", "")

	_, err := db.Conn.Exec(`
		INSERT INTO education (
			school,
			degree,
			field,
			startDate,
			endDate,
			grade,
			description
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		entry.School,
		entry.Degree,
		entry.Field,
		entry.StartDate,
		entry.EndDate,
		entry.Grade,
		entry.Description,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(ui.Success.Render("Education added successfully."))
}

func updateEducation(id int, field string) {
	if strings.EqualFold(field, "id") || !lib.HasField[models.Education](field) {
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
			"SELECT %s FROM education WHERE id = ? LIMIT 1",
			field,
		),
		id,
	).Scan(&currentValue)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(
				ui.Error.Render(
					"Education entry not found.",
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
			"UPDATE education SET %s = ? WHERE id = ?",
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
				"Education %d's %s updated successfully.",
				id,
				field,
			),
		),
	)
}

func deleteEducation(id int) {
	result, err := db.Conn.Exec(
		"DELETE FROM education WHERE id = ?",
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
				"Education entry not found.",
			),
		)
		return
	}

	fmt.Println(
		ui.Success.Render(
			fmt.Sprintf(
				"Education '%d' deleted successfully.",
				id,
			),
		),
	)
}

func education(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showEducation(false)
		return
	}

	if len(args) == 1 && args[0] == "--id" {
		showEducation(true)
		return
	}

	if len(args) == 1 && args[0] == "add" {
		addEducation()
		return
	}

	if len(args) == 3 && (args[1] == "--set" || args[1] == "-s") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid education id. Use a numeric id."))
			return
		}
		updateEducation(
			id,
			args[2],
		)
		return
	}

	if len(args) == 3 {
		if (args[1] == "--delete" || args[1] == "-d") && (args[2] == "--force" || args[2] == "-f") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid education id. Use a numeric id."))
				return
			}
			deleteEducation(id)
			return
		}

		if (args[1] == "--force" || args[1] == "-f") && (args[2] == "--delete" || args[2] == "-d") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid education id. Use a numeric id."))
				return
			}
			deleteEducation(id)
			return
		}
	}

	if len(args) == 2 && (args[1] == "--delete" || args[1] == "-d") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid education id. Use a numeric id."))
			return
		}

		fmt.Printf(
			"%s Are you sure you want to delete the education entry '%d'? This action cannot be undone. (y/N): ",
			ui.Warning.Render("Warning:"),
			id,
		)

		var confirmation string
		fmt.Scanln(&confirmation)

		if strings.ToLower(confirmation) != "y" {
			fmt.Println(ui.Info.Render("Education deletion cancelled."))
			return
		}
		deleteEducation(
			id,
		)
		return
	}

	fmt.Println(ui.Error.Render("Invalid command."))
}

var educationCmd = &cobra.Command{
	Use:                "education",
	Aliases:            []string{"edu"},
	Short:              "List and manage your education",
	Run:                education,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(educationCmd)
}
