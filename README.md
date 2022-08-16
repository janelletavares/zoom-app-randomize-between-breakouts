# Use Case

We run a social event over Webex with about 10 breakout sessions of two people each. The goal is that the same two people will not be paired up more than once.

# Behaviors

- authenticate with OAuth2
- upon start, ask me if the user wants to resume any "sessions" in progress
- fetch a list of participants
- randomize list of participants
- pair up participants and avoid matching the host, and if there is an odd number of users also avoid matching one cohost
- do not pair up the same participants more than once, if it can be avoided (A+B === B+A)
- write all sets of pairs to disk after each set of breakouts is opened
- match all people who have *** in front of their display name
BONUS:
- prompt the host to approve a set of matches
- shuffle one or more matches from a set before approving
- try to rotate the cohost that is left out of breakouts
- take a set of settings for the session (closed captioning, etc)
- take a set of settings for breakouts rooms

# Thoughts
does the user's unique id change between sessions?

# Command Line Usages

Build or refresh the binary
`make build`

## Get a new Access Token

export REDIRECT_URI=""
export CLIENT_ID=""
export CLIENT_SECRET=""
./webex-breakouts token

## List Meeting IDs
./webex-breakouts list-meetings


## Start managing a meeting's breakouts
export MEETING_ID=""
./webex-breakouts run
