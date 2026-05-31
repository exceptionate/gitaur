// cmd/version.go
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/exceptionate/gitaur/types"
	"github.com/spf13/cobra"
)

func getCommitHash() string {
	output, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return string(output)
}

var Version = ""
var Release = ""
var BuildDate = ""
var Commit = ""
var Platform = ""

func version(cmd *cobra.Command, args []string) {
	ver := types.Ver{
		Version:   Version,
		Release:   Release,
		BuildDate: BuildDate,
		Commit:    Commit,
		Platform:  Platform,
	}

	fmt.Printf("%+v\n", ver)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
