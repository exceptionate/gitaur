package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/exceptionate/gitaur/internal/auth"
	"github.com/exceptionate/gitaur/internal/db"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var force bool

func prompt(reader *bufio.Reader, label string) string {
	fmt.Print(ui.Label.Render(label + ": "))
	value, _ := reader.ReadString('\n')
	return strings.TrimSpace(value)
}

func setup(cmd *cobra.Command, args []string) {

	if err := db.Init(); err != nil {
		panic(err)
	}

	defer db.Conn.Close()

	if err := db.CreateSchema(); err != nil {
		panic(err)
	}

	userExists, err := db.UserExists()
	if err != nil {
		panic(err)
	}

	if userExists && !force {
		fmt.Println(
			ui.Warning.Render(
				"User profile and database already initialized. Use --force to overwrite.\n",
			),
		)
		return
	}

	if force {
		_, err = db.Conn.Exec("DELETE FROM user")
		if err != nil {
			panic(err)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(ui.Title.Render("Gitaur Setup"))
	fmt.Println()

	auth.Authenticate()

	token, err := keyring.Get(
		"gitaur",
		"github",
	)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	req.Header.Set(
		"Accept",
		"application/vnd.github+json",
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var githubUser struct {
		Name      string `json:"name"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	err = json.NewDecoder(resp.Body).Decode(
		&githubUser,
	)

	if err != nil {
		panic(err)
	}

	nameParts := strings.Fields(githubUser.Name)

	firstName := ""
	lastName := ""

	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}

	if len(nameParts) > 1 {
		lastName = strings.Join(
			nameParts[1:],
			" ",
		)
	}

	username := githubUser.Login
	email := githubUser.Email
	avatar := githubUser.AvatarURL

	linkedin := prompt(reader, "LinkedIn")
	portfolio := prompt(reader, "Portfolio")
	leetcode := prompt(reader, "LeetCode")
	phoneNumber := prompt(reader, "Phone Number")

	_, err = db.Conn.Exec(`
	INSERT INTO user (
		firstName,
		lastName,
		username,
		email,
		avatar,
		linkedin,
		portfolio,
		leetcode,
		phoneNumber
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		firstName,
		lastName,
		username,
		email,
		avatar,
		linkedin,
		portfolio,
		leetcode,
		phoneNumber,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(ui.Success.Render("Setup completed."))
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize Gitaur",
	Run:   setup,
}

func init() {
	setupCmd.Flags().BoolVarP(
		&force,
		"force",
		"f",
		false,
		"Replace existing profile",
	)

	rootCmd.AddCommand(setupCmd)
}
