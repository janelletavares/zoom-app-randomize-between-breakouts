package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/matching"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/io"
	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/webex"
)

func main() {
	ctx := context.TODO()

	args := os.Args[1:]

	switch command := args[0]; command {
	case "token":
		envVars, err := io.ReadConfigFromEnvironment()
		if err != nil {
			fmt.Printf("could not get credentials from environment: %v", err)
			os.Exit(1)
		}
		runOauthRedirectServer(envVars.ClientID, envVars.ClientSecret, envVars.RedirectURI, envVars.AccessToken)
	case "list-meetings":
		// grab token from file
		d, err := io.ReadDetailsFromFile()
		wbx := webex.New(ctx, d.AccessToken)
		response, err := wbx.ListMeetings()
		if err != nil {
			fmt.Printf("could not list meetings: %v", err)
			os.Exit(1)
		}
		fmt.Printf("MEETINGS:\n%s\n", spew.Sdump(response.Items))

	case "participants":
		// grab token from file
		d, err := io.ReadDetailsFromFile()
		wbx := webex.New(ctx, d.AccessToken)
		envVars, err := io.ReadConfigFromEnvironment()
		if err != nil {
			fmt.Printf("could not get meeting ID from environment: %v", err)
			os.Exit(1)
		}
		response, err := wbx.ListParticipants(envVars.MeetingID)
		if err != nil {
			fmt.Printf("could not list participants: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Participants:\n%s\n", spew.Sdump(response.Items))

	case "breakouts":
		// grab token from file
		d, err := io.ReadDetailsFromFile()
		wbx := webex.New(ctx, d.AccessToken)
		envVars, err := io.ReadConfigFromEnvironment()
		if err != nil {
			fmt.Printf("could not get meeting ID from environment: %v", err)
			os.Exit(1)
		}
		response, err := wbx.GetBreakoutSession(envVars.MeetingID)
		if err != nil {
			fmt.Printf("could not list breakouts: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Breakouts:\n%s\n", spew.Sdump(response.Items))

	case "run":
		var err error
		envVars, err := io.ReadConfigFromEnvironment()
		if err != nil {
			fmt.Printf("could not get meeting ID from environment: %v", err)
			os.Exit(1)
		}
		meetingID := envVars.MeetingID
		d, err := io.ReadDetailsFromFile()
		api := webex.New(ctx, d.AccessToken)
		var s *io.Session
		breakouts, err := s.SessionInProgress(meetingID)
		if err != nil {
			fmt.Printf("could not determine if there was a meeting in progress: %v\n", err)
		}
		if s != nil {
			fmt.Printf("Found session from %v. Resume this session (Y/n)?\n", breakouts.Timestamp)
			var response string
			_, _ = fmt.Scanln(&response) // wait for input to do it again
			if response == "Y" {
				matching.RegenerateMatchesHash(breakouts)
			}
		} else {
			s = &io.Session{}
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			}
			// cancel any pending processes with ctx
		}()
		for true {
			var response string
			fmt.Printf("Generating breakout sessions... Resume or end (R/e)? ")
			_, _ = fmt.Scanln(&response) // wait for input to do it again
			if response == "e" {
				err = io.WriteSessionToFile(s)
				if err != nil {
					fmt.Printf("trouble writing session to file: %v\n", err)
				}
				os.Exit(1)
			}

			// Generate matches
			breakout, err := matching.GenerateMatches(meetingID, api)
			if err != nil {
				fmt.Printf("OH NO! could not generate matches: %v\n", err)
			}
			fmt.Printf("Generated Matches:\n%s\n", spew.Sdump(breakout))
			// approve matches
			fmt.Printf("Approved (Y/n)? ")
			_, _ = fmt.Scanln(&response) // wait for input to do it again

			// start breakouts
			fmt.Printf("Trying to start breakouts...\n")
			res, err := api.PutBreakoutSession(meetingID, convertBreakoutToRequest(breakout))
			if err == nil && len(res.Items) == len(breakout) {
				fmt.Printf("Success!\n")
				fmt.Printf("Response: %s\n", spew.Sdump(res))

				// write latest breakout to file
				err = s.UpdateSession(meetingID, breakout)
				if err != nil {
					fmt.Printf("OH NO! could not update session to file: %v\n", err)
				}
			} else {
				fmt.Printf(
					"not successful: requested %d and received %d\n%s\n",
					len(breakout),
					len(res.Items),
					spew.Sdump(res))
				if err != nil {
					fmt.Printf("error: %s\n", err.Error())
				}
			}
		}

	default:
		fmt.Printf("%s is not a valid argument", command)
		os.Exit(1)
	}

}

func convertBreakoutToRequest(b matching.Breakout) webex.Matches {
	matches := make([]webex.Match, 0)

	for _, i := range b {
		webexMatch := webex.Match{}
		for _, j := range i.Participants {
			webexMatch = append(webexMatch, j.Email)
		}
		matches = append(matches, webexMatch)
	}
	return matches
}

func runOauthRedirectServer(clientID, clientSecret, redirectURI, accessToken string) {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{}

	// Create a new redirect route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the oauth endpoint
		// to get our access token

		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("client_id", clientID)
		data.Add("client_secret", clientSecret)
		data.Add("code", code)
		data.Add("redirect_uri", redirectURI)
		encodedData := data.Encode()

		req, err := http.NewRequest(http.MethodPost, "https://webexapis.com/v1/access_token", strings.NewReader(encodedData))
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		// We set this header since we want the response
		// as JSON
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Add("accept", "application/json")

		fmt.Printf("sending request to https://webexapis.com/v1/access_token\n")
		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		// Parse the request body into the `OAuthAccessResponse` struct
		var t OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		err = io.WriteDetailsToFile(clientID, clientSecret, redirectURI, accessToken)
		if err != nil {
			fmt.Printf("warning: issue writing details to file: %v", err)
		}

		// Finally, send a response to redirect the user to the "welcome" page
		// with the access token
		w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	fmt.Printf("starting to listen and serve....\n\n")
	http.ListenAndServe(":8080", nil)
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
