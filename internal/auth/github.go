package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/exceptionate/gitaur/internal/lib"
	"github.com/exceptionate/gitaur/internal/ui"
	"github.com/zalando/go-keyring"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
}

func Authenticate() {
	fmt.Println(ui.Info.Render("Perform Github authentication"))

	data, err := getDeviceCode()

	if err != nil {
		msg := lib.GetError(err)
		fmt.Println(ui.Error.Render(msg))
		return
	}

	fmt.Println()
	fmt.Println(ui.Info.Render("Open: " + data.VerificationURI))
	fmt.Println(ui.Info.Render("Code: " + data.UserCode))
	fmt.Println()

	err = waitForAuthentication(data)

	if err != nil {
		msg := lib.GetError(err)
		fmt.Println(ui.Error.Render(msg))
		return
	}

	fmt.Println(ui.Success.Render("Github authentication successful"))
}

func waitForAuthentication(data DeviceCodeResponse) error {
	for {
		time.Sleep(
			time.Duration(data.Interval) * time.Second,
		)

		body := bytes.NewBufferString(
			"client_id=Ov23li4OXflwPUZWQyWK" +
				"&device_code=" + data.DeviceCode +
				"&grant_type=urn:ietf:params:oauth:grant-type:device_code",
		)

		req, err := http.NewRequest(
			"POST",
			"https://github.com/login/oauth/access_token",
			body,
		)

		if err != nil {
			return err
		}

		req.Header.Set(
			"Accept",
			"application/json",
		)

		req.Header.Set(
			"Content-Type",
			"application/x-www-form-urlencoded",
		)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		var result map[string]any

		err = json.NewDecoder(resp.Body).Decode(
			&result,
		)

		resp.Body.Close()

		if err != nil {
			return err
		}

		token, ok := result["access_token"].(string)

		if ok {
			return keyring.Set(
				"gitaur",
				"github",
				token,
			)
		}

		e, ok := result["error"].(string)

		if !ok {
			continue
		}

		if e == "authorization_pending" {
			continue
		}

		return fmt.Errorf("%s", e)
	}
}

func getDeviceCode() (
	DeviceCodeResponse,
	error,
) {

	body := bytes.NewBufferString(
		"client_id=Ov23li4OXflwPUZWQyWK&scope=repo read:user",
	)

	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/device/code",
		body,
	)

	if err != nil {
		return DeviceCodeResponse{}, err
	}

	req.Header.Set(
		"Accept",
		"application/json",
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return DeviceCodeResponse{}, err
	}

	defer resp.Body.Close()

	var data DeviceCodeResponse

	err = json.NewDecoder(resp.Body).Decode(
		&data,
	)

	if err != nil {
		return DeviceCodeResponse{}, err
	}

	return data, nil
}
