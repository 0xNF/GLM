package notify

import (
	"bytes"
	"errors"
	"encoding/json"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"time"
	conf "github.com/0xNF/glm/src/internal/conf"
)


// validate checks that the email function can be properly executed
func validate(token string, channel string, msg string) error {
	if len(token) == 0 {
		return errors.New("Slack Token must not be empty")
	}
	if len(channel) == 0 {
		return errors.New("Slack Channel must not be empty")
	}
	if len(msg) == 0 {
		return errors.New("Slack message must not be empty")
	}
	return nil
}


// SendEmail sends a notification using the specified settings.
// Returns an error if things go wrong.
func SendEmail(emailConf *conf.GLMEmail, msg string) error {
	if err := validate(slackConf.Token, slackConf.Channel, msg); err != nil {
		return err
	}

	sl := &slack{
		Token: slackConf.Token,
		Channel: slackConf.Channel,
		Text: msg,
		Timeout: 30,
	}
	jobj, err := json.Marshal(sl)
	if err != nil {
		return err
	}
	os.Stdout.Write(jobj)
	req, err := http.NewRequest("POST", slackApIURL, bytes.NewBuffer(jobj))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+slackConf.Token)

	client := &http.Client{
		Timeout: time.Duration(time.Duration(30) * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	resBody := string(bodyBytes)
	if res.StatusCode != http.StatusOK {
		return errors.New(string(res.StatusCode))
	}
	if strings.Contains(resBody, "\"ok\":false") {
		return errors.New(resBody)
	}
	defer res.Body.Close()

	return nil
}