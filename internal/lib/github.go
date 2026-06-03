package lib

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/zalando/go-keyring"
)

func GetGithubToken() string {
	token, err := keyring.Get("gitaur", "github")
	if err != nil {
		panic(err)
	}

	return token
}

func GithubRequest(url string) (*http.Response, error) {
	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+GetGithubToken(),
	)

	return http.DefaultClient.Do(req)
}

func DecodeGithubContent(content string) string {
	decoded, _ := base64.StdEncoding.DecodeString(
		strings.ReplaceAll(content, "\n", ""),
	)

	return string(decoded)
}
