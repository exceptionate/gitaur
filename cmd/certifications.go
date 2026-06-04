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

func showCertifications(short bool) {
	rows, err := db.Conn.Query(`
		SELECT id, title, platform, description, url, date, tags
		FROM certifications
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var certifications models.Certifications

	for rows.Next() {
		var certification models.Certification
		var tagsJSON string

		err := rows.Scan(
			&certification.ID,
			&certification.Title,
			&certification.Platform,
			&certification.Description,
			&certification.Url,
			&certification.Date,
			&tagsJSON,
		)
		if err != nil {
			panic(err)
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &certification.Tags); err != nil {
				certification.Tags = nil
			}
		}

		certifications = append(certifications, certification)
	}

	if len(certifications) == 0 {
		fmt.Println(ui.Info.Render("No certifications found."))
		return
	}

	fmt.Printf("%s\n\n", ui.Info.Render(fmt.Sprintf("%d certifications found:", len(certifications))))

	if short {
		for _, certification := range certifications {
			fmt.Printf("%s: %s\n", ui.Label.Render(strconv.Itoa(certification.ID)), certification.Title)
		}
	} else {
		for i, certification := range certifications {
			certInfo := fmt.Sprintf(
				"%s %d\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
				ui.Label.Render("ID:"), certification.ID,
				ui.Label.Render("Title:"), certification.Title,
				ui.Label.Render("Platform:"), certification.Platform,
				ui.Label.Render("Description:"), certification.Description,
				ui.Label.Render("URL:"), certification.Url,
				ui.Label.Render("Date:"), certification.Date,
				ui.Label.Render("Tags:"), strings.Join(certification.Tags, ", "),
			)

			fmt.Println(certInfo)

			if i != len(certifications)-1 {
				fmt.Println("----------------------------")
			}
		}
	}
}

func addCertification() {
	fmt.Println(ui.Info.Render("Add Certification"))

	var certification models.Certification

	certification.Title = lib.PromptWithDefault("Title", "")
	certification.Platform = lib.PromptWithDefault("Platform", "")
	certification.Description = lib.PromptWithDefault("Description", "")
	certification.Url = lib.PromptWithDefault("URL", "")
	certification.Date = lib.PromptWithDefault("Date", "")
	tagsInput := lib.PromptWithDefault("Tags (comma separated)", "")
	if tagsInput != "" {
		certification.Tags = strings.Split(tagsInput, ",")
	}

	tagsJSON, err := json.Marshal(certification.Tags)
	if err != nil {
		panic(err)
	}

	_, err = db.Conn.Exec(`
		INSERT INTO certifications (
			title,
			platform,
			description,
			url,
			date,
			tags
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		certification.Title,
		certification.Platform,
		certification.Description,
		certification.Url,
		certification.Date,
		string(tagsJSON),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(ui.Success.Render("Certification added successfully."))
}

func updateCertification(id int, field string) {
	if strings.EqualFold(field, "id") || !lib.HasField[models.Certification](field) {
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
			"SELECT %s FROM certifications WHERE id = ? LIMIT 1",
			field,
		),
		id,
	).Scan(&currentValue)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(
				ui.Error.Render(
					"Certification not found.",
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
			"UPDATE certifications SET %s = ? WHERE id = ?",
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
				"Certification %d's %s updated successfully.",
				id,
				field,
			),
		),
	)
}

func deleteCertification(id int) {
	result, err := db.Conn.Exec(
		"DELETE FROM certifications WHERE id = ?",
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
				"Certification not found.",
			),
		)
		return
	}

	fmt.Println(
		ui.Success.Render(
			fmt.Sprintf(
				"Certification '%d' deleted successfully.",
				id,
			),
		),
	)
}

func certifications(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showCertifications(false)
		return
	}

	if len(args) == 1 && args[0] == "--id" {
		showCertifications(true)
		return
	}

	if len(args) == 1 && args[0] == "add" {
		addCertification()
		return
	}

	if len(args) == 3 && (args[1] == "--set" || args[1] == "-s") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid certification id. Use a numeric id."))
			return
		}
		updateCertification(
			id,
			args[2],
		)
		return
	}

	if len(args) == 3 {
		if (args[1] == "--delete" || args[1] == "-d") && (args[2] == "--force" || args[2] == "-f") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid certification id. Use a numeric id."))
				return
			}
			deleteCertification(id)
			return
		}

		if (args[1] == "--force" || args[1] == "-f") && (args[2] == "--delete" || args[2] == "-d") {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(ui.Error.Render("Invalid certification id. Use a numeric id."))
				return
			}
			deleteCertification(id)
			return
		}
	}

	if len(args) == 2 && (args[1] == "--delete" || args[1] == "-d") {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(ui.Error.Render("Invalid certification id. Use a numeric id."))
			return
		}

		fmt.Printf(
			"%s Are you sure you want to delete the certification '%d'? This action cannot be undone. (y/N): ",
			ui.Warning.Render("Warning:"),
			id,
		)

		var confirmation string
		fmt.Scanln(&confirmation)

		if strings.ToLower(confirmation) != "y" {
			fmt.Println(ui.Info.Render("Certification deletion cancelled."))
			return
		}
		deleteCertification(
			id,
		)
		return
	}

	fmt.Println(ui.Error.Render("Invalid command."))
}

var certificationsCmd = &cobra.Command{
	Use:                "certifications",
	Aliases:            []string{"cert", "certs"},
	Short:              "List and manage your certifications",
	Run:                certifications,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(certificationsCmd)
}
