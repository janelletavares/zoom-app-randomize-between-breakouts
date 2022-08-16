package io

import (
	"encoding/json"
	"os"

	"github.com/joeshaw/envdecode"

	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/matching"
)

const detailsFilename = ".details"
const sessionFilename = ".session"

type EnvVars struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURI  string `env:"REDIRECT_URI"`
	AccessToken  string `env:"ACCESS_TOKEN"`
	MeetingID    string `env:"MEETING_ID"`
}

type Details struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	AccessToken  string `json:"access_token"`
}

type Session struct {
	PerMeeting map[int]matching.BreakoutMatches
}

func ReadDetailsFromFile() (*Details, error) {
	d := Details{}
	bytes, err := os.ReadFile(detailsFilename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &d)
	return &d, err
}

func WriteDetailsToFile(clientID, clientSecret, redirectURI, accessToken string) error {
	d := &Details{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		AccessToken:  accessToken,
	}
	bytes, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return os.WriteFile(detailsFilename, bytes, 0644)
}

func ReadSessionFromFile() (*Session, error) {
	s := Session{}
	bytes, err := os.ReadFile(sessionFilename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &s)
	return &s, err

}

func ReadConfigFromEnvironment() (*EnvVars, error) {
	e := EnvVars{}
	if err := envdecode.Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

