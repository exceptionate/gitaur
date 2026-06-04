package cmd

import (
	"fmt"

	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/lib"
	"github.com/exceptionate/gitaur/internal/models"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
)

func showProfile() {
	var user models.User

	err := db.Conn.QueryRow(
		"SELECT firstName, lastName, username, email, avatar, linkedin, portfolio, leetcode, phoneNumber FROM user LIMIT 1",
	).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.LinkedIn,
		&user.Portfolio,
		&user.LeetCode,
		&user.PhoneNumber,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(ui.Info.Render("User Profile"))

	profileInfo := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
		ui.Label.Render("Name:"), user.FirstName+" "+user.LastName,
		ui.Label.Render("Username:"), user.Username,
		ui.Label.Render("Email:"), user.Email,
		ui.Label.Render("LinkedIn:"), user.LinkedIn,
		ui.Label.Render("Portfolio:"), user.Portfolio,
		ui.Label.Render("LeetCode:"), user.LeetCode,
		ui.Label.Render("Phone Number:"), user.PhoneNumber,
	)
	fmt.Println(profileInfo)

}

func updateProfile(field string) {
	if !lib.HasField[models.User](field) {
		fmt.Println(ui.Error.Render("Invalid field name. Please enter a valid field to update."))
		return
	}

	var value string

	fmt.Println(ui.Info.Render(fmt.Sprintf("You are updating the %s field: ", field)))
	fmt.Print(ui.Label.Render(fmt.Sprintf("%s :", field)))
	fmt.Scanln(&value)

	_, err := db.Conn.Exec(
		fmt.Sprintf("UPDATE user SET %s = ?", field),
		value,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(ui.Success.Render(fmt.Sprintf("%s updated successfully.", field)))
}

func profile(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		// show profile
		showProfile()

	} else if len(args) >= 2 {
		if args[0] == "update" {
			for i := 1; i < len(args); i++ {
				updateProfile(args[i])
				fmt.Println()
			}
		} else {
			fmt.Println(ui.Error.Render("Invalid subcommand."))
		}

	} else {
		fmt.Println(ui.Error.Render("Invalid command format."))
	}
}

func profileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "profile",
		Short: "View and manage user profile",
		Run:   profile,
	}
}

func init() {
	rootCmd.AddCommand(profileCmd())
}
