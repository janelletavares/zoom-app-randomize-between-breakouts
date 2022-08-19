package io

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"time"

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
	PerMeeting map[string]*matching.Breakouts
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

func WriteSessionToFile(s *Session) error {
	fmt.Printf("About to write Session to file...\n%s\n", spew.Sdump(s))
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(sessionFilename, bytes, 0644)
}

func ReadConfigFromEnvironment() (*EnvVars, error) {
	e := EnvVars{}
	if err := envdecode.Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *Session) SessionInProgress(meetingID string) (*matching.Breakouts, error) {
	if s == nil {
		var err error
		s, err = ReadSessionFromFile()
		if err != nil {
			return nil, err
		}
	}
	return s.PerMeeting[meetingID], nil
}

func (s *Session) UpdateSession(meetingID string, latest matching.Breakout) error {
	b := s.PerMeeting[meetingID]
	if b == nil {
		b = &matching.Breakouts{}
	}
	b.Timestamp = time.Now()
	b.Breakouts = append(b.Breakouts, latest)
	s.PerMeeting[meetingID] = b
	return WriteSessionToFile(s)
}
