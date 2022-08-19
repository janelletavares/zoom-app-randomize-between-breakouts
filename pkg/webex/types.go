package webex

type UserRole string

const (
	ParticipantRole UserRole = "PARTICIPANT"
	GuestRole       UserRole = "GUEST"
	HostRole        UserRole = "HOST"
	CohostRole      UserRole = "COHOST"
	PanelistRole    UserRole = "PANELIST"
	InterpreterRole UserRole = "INTERPRETER"
	PresenterRole   UserRole = "PRESENTER"
)

type MeetingType string

const (
	Meeting MeetingType = "MEETING"
	Event   MeetingType = "EVENT"
)

/*
   {
       "id": "870f51ff287b41be84648412901e0402_20191101T120000Z",
       "meetingSeriesId": "870f51ff287b41be84648412901e0402",
       "meetingNumber": "123456789",
       "title": "Example Daily Meeting",
       "agenda": "Example Agenda",
       "password": "BgJep@43",
       "phoneAndVideoSystemPassword": "12345678",
       "meetingType": "scheduledMeeting",
       "state": "ready",
       "isModified": false,
       "timezone": "UTC",
       "start": "2019-11-01T12:00:00Z",
       "end": "2019-11-01T13:00:00Z",
       "hostUserId": "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9jN2ZkNzNmMi05ZjFlLTQ3ZjctYWEwNS05ZWI5OGJiNjljYzY",
       "hostDisplayName": "John Andersen",
       "hostEmail": "john.andersen@example.com",
       "hostKey": "123456",
       "siteUrl": "site4-example.webex.com",
       "webLink": "https://site4-example.webex.com/site4/j.php?MTID=md41817da6a55b0925530cb88b3577b1e",
       "sipAddress": "123456789@site4-example.webex.com",
       "dialInIpAddress": "192.168.100.100",
       "enabledAutoRecordMeeting": false,
       "allowAnyUserToBeCoHost": false,
       "enabledJoinBeforeHost": false,
       "enableConnectAudioBeforeHost": false,
       "joinBeforeHostMinutes": 0,
       "excludePassword": false,
       "publicMeeting": false,
       "reminderTime": 10,
       "unlockedMeetingJoinSecurity": "allowJoin",
       "sessionTypeId": 3,
       "enableAutomaticLock": false,
       "automaticLockMinutes": 0,
       "allowFirstUserToBeCoHost": false,
       "allowAuthenticatedDevices": false,
       "telephony": {
           "accessCode": "1234567890",
           "callInNumbers": [
               {
                   "label": "US Toll",
                   "callInNumber": "123456789",
                   "tollType": "toll"
               }
           ],
           "links": [
               {
                   "rel": "globalCallinNumbers",
                   "href": "/api/v1/meetings/870f51ff287b41be84648412901e0402_20191101T120000Z/globalCallinNumbers",
                   "method": "GET"
               }
           ]
       },
       "meetingOptions": {
           "enabledChat": true,
           "enabledVideo": true,
           "enabledPolling": false,
           "enabledNote": true,
           "noteType": "allowAll",
           "enabledClosedCaptions": false,
           "enabledFileTransfer": false,
           "enabledUCFRichMedia": false
       },
       "integrationTags": [
           "dbaeceebea5c4a63ac9d5ef1edfe36b9",
           "85e1d6319aa94c0583a6891280e3437d",
           "27226d1311b947f3a68d6bdf8e4e19a1"
       ],
       "scheduledType": "meeting",
       "enabledBreakoutSessions": true,
       "links": [
           {
               "rel": "breakoutSessions",
               "href":"/v1/meetings/870f51ff287b41be84648412901e0402/breakoutSessions",
               "method": "GET"
           }
       ]
   },

*/

type MeetingOptions struct {
	EnabledChat           bool   `json:"enabledChat"`
	EnabledVideo          bool   `json:"enabledVideo"`
	EnabledPolling        bool   `json:"enabledPolling"`
	EnabledNote           bool   `json:"enabledNote"`
	NoteType              string `json:"noteType"`
	EnabledClosedCaptions bool   `json:"enabledClosedCaptions"`
	EnabledFileTransfer   bool   `json:"enabledFileTransfer"`
	EnabledUCFRichMedia   bool   `json:"enabledUCFRichMedia"`
}

type MeetingResponse struct {
	Items []MeetingDetails `json:"items"`
}

type MeetingDetails struct {
	ID              string `json:"id"` //-- Meeting ID. If app.isPrivateDataAvailable is true the value is a real meeting ID, otherwise it's a derived meeting ID. Derived IDs are guaranteed to be consistent for all users of the meeting.
	MeetingNumber   string `json:"meetingNumber"`
	State           string `json:"state"`
	HostDisplayName string `json:"hostDisplayName"`
	Title           string `json:"title"` //-- Title of the given meeting; blank if app.isPrivateDataAvailable is false.
	StartTime       string `json:"start"` //-- Start time for a scheduled meeting in UTC ("2021-01-17T13:00:00.00Z", for example). Blank for personal meeting rooms.
	EndTime         string `json:"end"`   //-- End time for a scheduled meeting in UTC ("2021-01-17T13:00:10.00Z", for example). Blank for personal meeting rooms.
	//	ConferenceId            string         `json:"conferenceId"` //-- Conference ID. A unique ID that's created when then the first participant joins a meeting. If app.isPrivateDataAvailable is true then the value is a real conference ID, otherwise it's a derived conference ID. Derived IDs are guaranteed to be consistent for all users of the meeting.
	//	UserRoles               []UserRole     `json:"userRoles"`
	Type                    MeetingType    `json:"scheduledType"`
	EnabledBreakoutSessions bool           `json:"enabledBreakoutSessions"`
	Options                 MeetingOptions `json:"meetingOptions"`
}

/*
{
    "hostEmail": "john.andersen@example.com",
    "sendEmail": true,
    "items": [
        {
            "name": "Breakout Session 1",
            "invitees": [
                "rachel.green@example.com",
                "monica.geller@example.com"
            ]
        },
        {
            "name": "Breakout Session N",
            "invitees": [
                "ross.geller@example.com",
                "chandler.bing@example.com"
            ]
        }
    ]
}

*/

/*
   {
     "registrantId": "123456",
     "status": "pending",
     "firstName": "bob",
     "lastName": "Lee",
     "email": "bob@cisco.com",
     "jobTitle": "manager",
     "companyName": "cisco",
     "address1": "address1 string",
     "address2": "address2 string",
     "city": "New York",
     "state": "New York",
     "zipCode": 123456,
     "countryRegion": "America",
     "workPhone": "+1 123456",
     "fax": "123456",
     "registrationTime": "2021-09-07T09:29:13+08:00"
*/
type Participant struct {
	ID          string `json:"id"`
	Host        bool   `json:"host"`
	CoHost      bool   `json:"coHost"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	State       string `json:"state"`
	JoinedTime  string `json:"joinedTime"`
}

type ParticipantsResponse struct {
	Items []Participant `json:"items"`
}

type BreakoutSessionsRequest struct {
	HostEmail string     `json:"hostEmail"`
	SendEmail bool       `json:"sendEmail"`
	Items     []Breakout `json:"items"`
}

type BreakoutSessionsResponse struct {
	Items []Breakout `json:"items"`
}

type Breakout struct {
	Name     string   `json:"name"`
	Invitees []string `json:"invitees"`
}

type Match []string // list of emails
type Matches []Match
