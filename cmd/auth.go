package cmd

import (
	"fmt"

	"github.com/exceptionate/gitaur/internal/auth"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
)

var refreshAuth bool

func authCommand(cmd *cobra.Command, args []string) {
	if !refreshAuth {
		fmt.Println(ui.Warning.Render("Use --refresh to re-authenticate."))
		return
	}

	auth.Authenticate()
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Github authentication",
	Run:   authCommand,
}

func init() {
	authCmd.Flags().BoolVar(
		&refreshAuth,
		"refresh",
		false,
		"Redo login and update token",
	)

	rootCmd.AddCommand(authCmd)
}
