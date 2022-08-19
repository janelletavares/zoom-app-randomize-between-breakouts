package matching

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/janelletavares/webex-app-randomize-between-breakouts/pkg/webex"
	"strings"
	"time"
)

var matchesHash map[string]bool

const DoNotMatch = "***"

type Breakouts struct {
	Timestamp time.Time // last written to file
	Breakouts []Breakout
}

type Breakout []Match

type Match struct {
	Participants []Participant // usually 2, expect when there's ***
	Duplicate    bool
}

type Participant struct {
	ID          string
	DisplayName string
	Email       string
}

func GenerateMatches(meetingID string, api webex.API) (Breakout, error) {
	if matchesHash == nil {
		matchesHash = make(map[string]bool)
	}
	particpantsResponse, err := api.ListParticipants(meetingID)
	if err != nil {
		return nil, fmt.Errorf("unable to list participants: %v", err)
	}
	participants := particpantsResponse.Items
	fmt.Printf("Found %d participants\n", len(participants))
	// @TODO sort

	// put all *** in one room
	matches := make([]Match, 0)
	i := 0
	m := Match{Participants: make([]Participant, 0)}
	for strings.Contains(participants[i].DisplayName, DoNotMatch) {
		m.Participants = append(
			m.Participants,
			Participant{DisplayName: participants[i].DisplayName, Email: participants[i].Email})
		i++
	}
	if i != 0 {
		fmt.Printf("Sequestered %d folks who did not want to be matched\n", i)
		matches = append(matches, m)
	}

	// match and check
	//len := len(participants) - i
	fmt.Printf("Remainging participants to match: %d\n", len(participants[i:]))
	for len(participants[i:]) > 1 {
		second := i + 1
		dup := Duplicate(participants[i].ID, participants[second].ID)
		tryingToResolve := false
		for dup && second < (len(participants)-2) {
			tryingToResolve = true
			fmt.Printf("found a duplicate and trying to fix\n")
			second++
			dup = Duplicate(participants[i].ID, participants[second].ID)
		}

		if !dup {
			RegisterMatch(participants[i].ID, participants[second].ID)
		}
		if tryingToResolve && dup {
			fmt.Printf("Unable to avoid duplicate!\n")
		}
		fmt.Printf("Matching folks at index %d and %d\n", i, second)
		m := Match{Participants: []Participant{
			{DisplayName: participants[i].DisplayName, Email: participants[i].Email},
			{DisplayName: participants[second].DisplayName, Email: participants[second].Email},
		},
			Duplicate: dup}
		matches = append(matches, m)
		participants = remove(participants, second)
		participants = remove(participants, i)
	}

	/*
		// extra 2 or 3 rooms for manual shenanigans
		fmt.Printf("Adding two more empty sandbox breakouts\n")
		for j := 0; j < 2; j++ {
			matches = append(matches, Match{})
		}

	*/
	return matches, nil
}

func remove(slice []webex.Participant, index int) []webex.Participant {
	return append(slice[:index], slice[index+1:]...)
}

func Duplicate(a, b string) bool {
	fmt.Printf("Assessing duplicate...\n%s\n", spew.Sdump(matchesHash))
	oka := matchesHash[a+b]
	okb := matchesHash[b+a]
	return oka == false && okb == false
}

func RegisterMatch(a, b string) {
	matchesHash[a+b] = true
	matchesHash[b+a] = true
}

func (b *Breakouts) UpdateMatch() {
	// @TODO handle a rejected breakout
}

func RegenerateMatchesHash(b *Breakouts) error {
	matchesHash = make(map[string]bool)
	for _, b := range b.Breakouts {
		for _, m := range b {
			if !strings.Contains(m.Participants[0].DisplayName, DoNotMatch) {
				RegisterMatch(m.Participants[0].ID, m.Participants[0].ID)
			}
		}
	}
	return nil
}
