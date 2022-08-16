package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/joeshaw/envdecode"

	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/io"
	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/webex"
)

// could be environment variables
var clientID = "C083e421d6d25ee2b724c759534f94ef94ddfb219894f0361e205a55913e5a71b"
var clientSecret = "80ad838015a71f72ce10bbb4c9ab11cace8e71254e2e4ada1e13b49a126291f2"
var redirectURI = "http://localhost:8080/oauth/redirect"

var accessToken = "OTE1MDYyNDQtOGIwYi00OGQ0LThiNjktNDk0YzljNzYyZmZjNWEwNjg5NWEtZDky_P0A1_efbefec4-802f-439d-b9f8-0704addd89c2"

func main() {
	ctx := context.TODO()

	args := os.Args[1:]

	switch command := args[0]; command {
	case "token":
		// grab values from environment
		if err := envdecode.Decode(&c); err != nil {
			fmt.Printf("could not get environment variables: %v", err)
			os.Exit(1)
		}
		//runOauthRedirectServer()
	case "list-meetings":
		// grab token from file
		d, err := io.ReadDetailsFromFile()
		wbx := webex.New(d.AccessToken)
		response, err := wbx.ListMeetings(ctx)
		if err != nil {
			fmt.Printf("could not list meetings: %v", err)
			os.Exit(1)
		}
		fmt.Printf("meetings:\n%s\n", spew.Sdump(response.Items))
	case "run":
		// check for session with the same meeting ID
		if err := envdecode.Decode(&c); err != nil {
			return nil, err
		}
		//  fmt.Scanln(&first)
	default:
		fmt.Printf("%s is not a valid argument", command)
		os.Exit(1)
	}

}

func runOauthRedirectServer() {
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
