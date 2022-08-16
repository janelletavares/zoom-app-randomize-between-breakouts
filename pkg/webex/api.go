package webex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
GET /rooms
Authorization: Bearer THE_ACCESS_TOKEN
Accept: application/json

*/

type WebexAPI interface {
	ListMeetings(ctx context.Context) (*MeetingResponse, error)
}

type api struct {
	accessToken string
}

func New(accessToken string) WebexAPI {
	return &api{accessToken: accessToken}
}

func (a *api) ListMeetings(ctx context.Context) (*MeetingResponse, error) {
	meetingResponse := &MeetingResponse{}
	res, err := a.requestWrapper(ctx, http.MethodGet, "https://webexapis.com/v1/meetings", nil)
	defer res.Body.Close()
	if err != nil {
		return meetingResponse, err
	}

	/*
		buf := new(strings.Builder)
		_, err = io.Copy(buf, res.Body)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		fmt.Println(buf.String())
	*/

	var response MeetingResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("could not parse JSON response: %v", err)
	}
	return meetingResponse, nil
}

// GET https://webexapis.com/v1/meetingParticipants
// GET  https://webexapis.com/v1/meetingPreferences
// PUT https://webexapis.com/v1/meetingPreferences/personalMeetingRoom
// PUT /v1/meetings/￼meetingId/breakoutSessions
// GET /v1/meetings/{meetingId}/breakoutSessions

func (a *api) requestWrapper(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {

	httpClient := http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("could not send HTTP request: %v", err)
	}
	if res.StatusCode > 204 {
		return res, fmt.Errorf("unsuccessful response: %d %s", res.StatusCode, res.Status)
	}
	return res, nil
}
