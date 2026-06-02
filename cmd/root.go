package cmd

import (
	"fmt"

	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
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

		//check if db user table has any entries, if not prompt user to run setup
		//db exists at this point, so we can safely open it
		if db.Conn == nil {
			err := db.Init()
			if err != nil {
				panic(err)
			}
		}

		var count int
		err := db.Conn.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
		if err != nil {
			panic(err)
		}

		if count == 0 && cmd.Name() != "setup" {
			fmt.Println(
				ui.Warning.Render(
					"\nNo user found in database. Run 'gitaur setup'\n",
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
