package webex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"net/http"
)

/*
GET /rooms
Authorization: Bearer THE_ACCESS_TOKEN
Accept: application/json

*/

type API interface {
	ListMeetings() (*MeetingResponse, error)
	ListParticipants(meetingID string) (*ParticipantsResponse, error)
	GetBreakoutSession(meetingID string) (*BreakoutSessionsResponse, error)
	DeleteBreakoutSessions(meetingID string) (*BreakoutSessionsResponse, error)
	PutBreakoutSession(meetingID string, breakout []Match) (*BreakoutSessionsResponse, error)
}

type api struct {
	ctx         context.Context
	accessToken string
}

func New(ctx context.Context, accessToken string) API {
	return &api{ctx: ctx, accessToken: accessToken}
}

func (a *api) ListMeetings() (*MeetingResponse, error) {
	meetingResponse := &MeetingResponse{}
	res, err := a.requestWrapper(a.ctx, http.MethodGet, "https://webexapis.com/v1/meetings", nil)
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

	if err := json.NewDecoder(res.Body).Decode(&meetingResponse); err != nil {
		return nil, fmt.Errorf("could not parse JSON response: %v", err)
	}
	return meetingResponse, nil
}

func (a *api) ListParticipants(meetingID string) (*ParticipantsResponse, error) {
	// GET https://webexapis.com/v1/meetingParticipants
	response := &ParticipantsResponse{}
	res, err := a.requestWrapper(a.ctx,
		http.MethodGet,
		fmt.Sprintf("https://webexapis.com/v1/meetingParticipants?meetingId=%s", meetingID),
		nil)
	defer res.Body.Close()
	if err != nil {
		return response, err
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("could not parse JSON response: %v", err)
	}
	return response, nil
}

func (a *api) PutBreakoutSession(meetingID string, breakout []Match) (*BreakoutSessionsResponse, error) {
	// PUT /v1/meetings/{meetingId}/breakoutSessions
	response := &BreakoutSessionsResponse{}
	body := &BreakoutSessionsRequest{
		HostEmail: "jtav77@gmail.com", // @TODO fix me
		SendEmail: false,
		Items:     make([]Breakout, 0),
	}
	for i, m := range breakout {
		b := Breakout{Name: fmt.Sprintf("Breakout Room %d", i)}
		b.Invitees = make([]string, 0)
		for _, p := range m {
			b.Invitees = append(b.Invitees, p)
		}

		body.Items = append(body.Items, b)
	}
	fmt.Printf("Request to start Breakout Session:\n%s\n", spew.Sdump(body))
	buf := new(bytes.Buffer)
	if body != nil {
		if err := encodeBody(buf, body); err != nil {
			return nil, err
		}
	}
	res, err := a.requestWrapper(a.ctx,
		http.MethodGet,
		fmt.Sprintf("https://webexapis.com/v1/meetings/%s/breakoutSessions", meetingID),
		buf)
	defer res.Body.Close()
	if err != nil {
		return response, err
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("could not parse JSON response: %v", err)
	}
	return response, nil
}

func encodeBody(w io.Writer, v interface{}) error {
	switch body := v.(type) {
	case string:
		_, err := w.Write([]byte(body))
		return err
	case []byte:
		_, err := w.Write(body)
		return err
	default:
		return json.NewEncoder(w).Encode(v)
	}
}

func (a *api) GetBreakoutSession(meetingID string) (*BreakoutSessionsResponse, error) {
	// GET /v1/meetings/{meetingId}/breakoutSessions
	response := &BreakoutSessionsResponse{}
	res, err := a.requestWrapper(a.ctx,
		http.MethodGet,
		fmt.Sprintf("https://webexapis.com/v1/meetings/%s/breakoutSessions", meetingID),
		nil)
	defer res.Body.Close()
	if err != nil {
		return response, err
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("could not parse JSON response: %v", err)
	}
	return response, nil
}

func (a *api) DeleteBreakoutSessions(meetingID string) (*BreakoutSessionsResponse, error) {
	// DELETE /v1/meetings/{meetingId}/breakoutSessions
	response := &BreakoutSessionsResponse{}
	res, err := a.requestWrapper(a.ctx,
		http.MethodDelete,
		fmt.Sprintf("https://webexapis.com/v1/meetings/%s/breakoutSessions", meetingID),
		nil)
	defer res.Body.Close()
	if err != nil {
		return response, err
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("could not parse JSON response: %v", err)
	}
	return response, nil
}

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
	fmt.Printf("Response code: %d\n", res.StatusCode)
	if res.StatusCode > 204 {
		return res, fmt.Errorf("unsuccessful response: %d %s", res.StatusCode, res.Status)
	}
	return res, nil
}
