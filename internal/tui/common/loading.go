package common

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// loadingMessages is a list of messages displayed while loading.
// They are copied from sonarr and radarr.
var loadingMessages = [...]string{
	"Downloading more RAM",
	"Now in Technicolor",
	"Previously on Sonarr...",
	"Previously on Radarr...",
	"Bleep Bloop.",
	"Locating the required gigapixels to render...",
	"Spinning up the hamster wheel...",
	"At least you're not on hold",
	"Hum something loud while others stare",
	"Loading humorous message... Please Wait",
	"I could've been faster in Python",
	"Don't forget to rewind your episodes",
	"Don't forget to rewind your movies",
	"Congratulations! you are the 1000th visitor.",
	"HELP! I'm being held hostage and forced to write these stupid lines!",
	"RE-calibrating the internet...",
	"I'll be here all week",
	"Don't forget to tip your waitress",
	"Apply directly to the forehead",
	"Loading Battlestation",
}

// GetRandomLoadingMessage returns a random loading message.
func GetRandomLoadingMessage() string {
	// #nosec G404
	return loadingMessages[rand.Intn(len(loadingMessages))]
}
