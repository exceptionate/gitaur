// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
)

var Version = "DEFAULT_VERSION"
var Release = "DEFAULT_RELEASE"
var BuildDate = "DEFAULT_DATE"
var Commit = "DEFAULT_COMMIT"
var Platform = "DEFAULT_PLATFORM"

func version(cmd *cobra.Command, args []string) {

	fmt.Println(ui.Title.Render("Gitaur CLI"))
	fmt.Println()

	fmt.Printf("%s %s\n", ui.Label.Render("Version:"), Version)
	fmt.Printf("%s %s\n", ui.Label.Render("Release:"), Release)
	fmt.Printf("%s %s\n", ui.Label.Render("Build:"), BuildDate)
	fmt.Printf("%s %s\n", ui.Label.Render("Commit:"), Commit[:7])
	fmt.Printf("%s %s\n", ui.Label.Render("Platform:"), Platform)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
