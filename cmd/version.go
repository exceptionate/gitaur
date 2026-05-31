// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "DEFAULT_VERSION"
var Release = "DEFAULT_RELEASE"
var BuildDate = "DEFAULT_DATE"
var Commit = "DEFAULT_COMMIT"
var Platform = "DEFAULT_PLATFORM"

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("Gitaur CLI %s (%s)\n\nBuilt:    %s\nCommit:   %s\nPlatform: %s\n",
		Version,
		Release,
		BuildDate,
		Commit[:7],
		Platform,
	)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
