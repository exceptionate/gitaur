package cmd

import (
	"fmt"

	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func preRun(cmd *cobra.Command, args []string) {
	if cmd.Name() != "setup" {
		if !db.Exists() {

			fmt.Println(
				ui.Warning.Render(
					"\nDatabase not initialized. Run 'gitaur setup'\n",
				),
			)
			return
		}

		//db exists at this point, so we can safely open it
		if db.Conn == nil {
			err := db.Init()
			if err != nil {
				panic(err)
			}
		}

		SchemaExists, err := db.SchemaExists()
		if err != nil {
			panic(err)
		}

		if !SchemaExists {
			fmt.Println(
				ui.Warning.Render(
					"\nDatabase schema is broken. Run 'gitaur setup' to recreate it.\n",
				),
			)
			return
		}

		//if user exists
		UserExists, err := db.UserExists()
		if err != nil {
			panic(err)
		}

		if !UserExists {
			fmt.Println(
				ui.Warning.Render(
					"\nNo user profile found. Run 'gitaur setup'\n",
				),
			)
			return
		}

		// Auth check. check from keyring if token exists, if not prompt user to run auth
		_, err = keyring.Get("gitaur", "github")
		if err != nil {
			fmt.Println(
				ui.Warning.Render(
					"\nNo Github token found. Run 'gitaur auth --refresh'\n",
				),
			)
			return
		}

	}
}

var rootCmd = &cobra.Command{
	Use:              "gitaur",
	Short:            "Gitaur CLI",
	PersistentPreRun: preRun,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
