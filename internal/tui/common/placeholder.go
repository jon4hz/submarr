package common

import "math/rand"

// some copied from other depths of the internet
// some created by github copilot (I'm not particularly proud of that but they dont seem to bad)
var placeholders = [...]string{
	"Land was created to provide a place for boats to visit",
	"A Pirate's favorite movie is one that is rated \"ARRRR\"!",
	"Arrrrrgh!",
	"Avast ye!",
	"Avast ye scurvy dogs!",
	"Avast! Pull Me Mast!",
	"Avast! Ye landlubbers!",
	"Ahoy! lets trouble the water!",
	"Why is the rum gone?",
	"Let's jump on board, and cut them to pieces",
	"Surrrrrender the booty!",
	"Shiver me timbers! Me wooden leg has termites",
	"Swab My Deck, Wench",
	"Keep calm and say \"Arrr\"",
	"Rubbers are for land lubbers",
	"Piracy is the way o life. Ahoy",
	"Pirates do it harrrrrder!",
	"Shut Ye Pie Hole, I'm Diving in Ye Bung Hole",
	"The Code is more like guidelines, really",
	"I be ruler of the seven seas!",
	"Prepare to be boarded.",
	"Yo ho ho and a bottle of rum.",
	"Thar she blows!",
	"To life, love and loot.",
	"It is a glorious thing to be a Pirate King.",
	"Not all treasure's silver and gold, mate.",
	"When a pirate grows rich enough, they make him a prince.",
	"Ah. Love. A dreadful bond.",
	"If rum can't fix it, you are not using enough rum.",
	"All for rum and rum for all.",
	"Money can't buy you happiness… but it can buy you rum!",
	"Sometimes it just takes a pirate to get the job done.",
	"Damn you villains, who are you? And from whence came you?",
}

func GetRandomPlaceholder() string {
	// #nosec G404
	return placeholders[rand.Intn(len(placeholders))]
}
