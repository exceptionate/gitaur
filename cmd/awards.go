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

func showAwards(short bool) {
	rows, err := db.Conn.Query(`
		SELECT id, title, issuer, type, tags, description, date
		FROM awards
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var awards models.Awards

	for rows.Next() {
		var award models.Award
		var tagsJSON string

		err := rows.Scan(
			&award.ID,
			&award.Title,
			&award.Issuer,
			&award.Type,
			&tagsJSON,
			&award.Description,
			&award.Date,
		)
		if err != nil {
			panic(err)
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &award.Tags); err != nil {
				award.Tags = nil
			}
		}

		awards = append(awards, award)
	}

	if len(awards) == 0 {
		fmt.Println(ui.Info.Render("No awards found."))
		return
	}

	fmt.Printf("%s\n\n", ui.Info.Render(fmt.Sprintf("%d awards found:", len(awards))))

	if short {
		for _, award := range awards {
			fmt.Printf("%s: %s\n", ui.Label.Render(strconv.Itoa(award.ID)), award.Title)
		}
	} else {
		for i, award := range awards {
			awardInfo := fmt.Sprintf(
				"%s %d\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
				ui.Label.Render("ID:"), award.ID,
				ui.Label.Render("Title:"), award.Title,
				ui.Label.Render("Issuer:"), award.Issuer,
				ui.Label.Render("Type:"), award.Type,
				ui.Label.Render("Tags:"), strings.Join(award.Tags, ", "),
				ui.Label.Render("Description:"), award.Description,
				ui.Label.Render("Date:"), award.Date,
			)

			fmt.Println(awardInfo)

			if i != len(awards)-1 {
				fmt.Println("----------------------------")
			}
		}
	}
}

func addAward() {
	fmt.Println(ui.Info.Render("Add Award"))

	var award models.Award

	award.Title = lib.PromptWithDefault("Title", "")
	award.Issuer = lib.PromptWithDefault("Issuer", "")
	award.Type = lib.PromptWithDefault("Type", "")
	tagsInput := lib.PromptWithDefault("Tags (comma separated)", "")
	if tagsInput != "" {
		award.Tags = strings.Split(tagsInput, ",")
	}
	award.Description = lib.PromptWithDefault("Description", "")
	award.Date = lib.PromptWithDefault("Date", "")

	tagsJSON, err := json.Marshal(award.Tags)
	if err != nil {
		panic(err)
	}

	_, err = db.Conn.Exec(`
		INSERT INTO awards (
			title,
			issuer,
			type,
			tags,
			description,
			date
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		award.Title,
		award.Issuer,
		award.Type,
		string(tagsJSON),
		award.Description,
		award.Date,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(ui.Success.Render("Award added successfully."))
}

func updateAward(id int, field string) {
	if strings.EqualFold(field, "id") || !lib.HasField[models.Award](field) {
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
			"SELECT %s FROM awards WHERE id = ? LIMIT 1",
			field,
		),
		id,
	).Scan(&currentValue)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(
				ui.Error.Render(
					"Award not found.",
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
			"UPDATE awards SET %s = ? WHERE id = ?",
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
				"Award %d's %s updated successfully.",
				id,
				field,
			),
		),
	)
}

func deleteAward(id int) {
	result, err := db.Conn.Exec(
		"DELETE FROM awards WHERE id = ?",
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
				"Award not found.",
			),
		)
		return
	}

	fmt.Println(
		ui.Success.Render(
			fmt.Sprintf(
				"Award '%d' deleted successfully.",
				id,
			),
		),
	)
}

func awards(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showAwards(false)
		return
	}

	if len(args) == 1 && args[0] == "--id" {
		showAwards(true)
		return
	}

	if len(args) == 1 && args[0] == "add" {
		addAward()
		return
	}

	if len(args) == 3 && (args[1] == "--set" || args[1] == "-s") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid award id. Use a numeric id."))
			return
		}
		updateAward(
			id,
			args[2],
		)
		return
	}

	if len(args) == 3 {
		if (args[1] == "--delete" || args[1] == "-d") && (args[2] == "--force" || args[2] == "-f") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid award id. Use a numeric id."))
				return
			}
			deleteAward(id)
			return
		}

		if (args[1] == "--force" || args[1] == "-f") && (args[2] == "--delete" || args[2] == "-d") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid award id. Use a numeric id."))
				return
			}
			deleteAward(id)
			return
		}
	}

	if len(args) == 2 && (args[1] == "--delete" || args[1] == "-d") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid award id. Use a numeric id."))
			return
		}

		fmt.Printf(
			"%s Are you sure you want to delete the award '%d'? This action cannot be undone. (y/N): ",
			ui.Warning.Render("Warning:"),
			id,
		)

		var confirmation string
		fmt.Scanln(&confirmation)

		if strings.ToLower(confirmation) != "y" {
			fmt.Println(ui.Info.Render("Award deletion cancelled."))
			return
		}
		deleteAward(
			id,
		)
		return
	}

	fmt.Println(ui.Error.Render("Invalid command."))
}

var awardsCmd = &cobra.Command{
	Use:                "awards",
	Short:              "List and manage your awards",
	Run:                awards,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(awardsCmd)
}
