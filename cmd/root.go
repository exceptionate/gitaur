package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gitaur",
	Short: "Gitaur CLI",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
