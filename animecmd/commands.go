package animecmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jasonpuglisi/ircutil"
)

// Init adds animecmd's functions to the command map.
func Init(cmdMap ircutil.CmdMap) {
	ircutil.AddCommand(cmdMap, "inami/animecmd.Search", Search)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Watch", Watch)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Countdown", Countdown)
}

// Search searches an anime database and returns relevant search results.
// Function key: inami/animecmd.Search
func Search(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Search Hummingbird for shows matching query.
	shows, err := hummingbirdSearch(strings.Join(message.Args, " "))
	if err != nil {
		log.Println(err)
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error fetching shows, try again later")
		return
	}

	// Send response if no shows are found.
	if len(shows) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			"No shows found")
		return
	}

	// Send response with found shows and their URLs.
	ircutil.SendResponse(client, message.Source, message.Target,
		"Shows found:")
	for _, s := range shows {
		ircutil.SendResponse(client, message.Source, message.Target,
			fmt.Sprintf("- %s (%s)", s.Title, s.URL))
	}
}

// Watch gets data of a show and starts a group watching session.
// Function key: inami/animecmd.Watch
func Watch(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Send response if episode number is not a number.
	numStr := message.Args[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Println(err)
		ircutil.SendResponse(client, message.Source, message.Target,
			"Invalid episode number")
		return
	}

	// Get show data from Hummingbird.
	show, err := hummingbirdShow(message.Args[0])
	if err != nil {
		log.Println(err)
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error fetching show, try again later")
		return
	}

	// Send response if show not found.
	if show.Anime.Slug == "" {
		ircutil.SendResponse(client, message.Source, message.Target, "Show not "+
			"found, make sure you're using the name after /anime/ in the URL")
		return
	}

	// Split anime and episode information, and find specific episode.
	anime, episodes := show.Anime, show.Linked.Episodes
	var episode *episode
	for i := range episodes {
		e := &episodes[i]
		if e.Number == num {
			episode = e
		}
	}

	// Clear episode title if it doesn't exist.
	episodeTitle := ""
	if episode != nil && episode.Title != "Episode "+numStr {
		episodeTitle = " \"" + episode.Title + "\""
	}

	// Start countdown and sleep before sending episode title.
	Countdown(client, command, message)
	time.Sleep(time.Second * 5)

	// Send response with episode information.
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("You're watching %s Episode %s%s", anime.Titles.Canonical,
			numStr, episodeTitle))
}

// Countdown starts a countdown to coordinate group watching.
// Function key: inami.animeutil/Countdown
func Countdown(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Send response for countdown start.
	ircutil.SendResponse(client, message.Source, message.Target,
		"Starting countdown, press play when I say \"Start!\"")

	// Send response with seconds remaining, or "Start!" at 0, and decrement
	// seconds remaining.
	ticker := time.NewTicker(time.Second)
	i := 6
	for range ticker.C {
		s := strconv.Itoa(i)
		if i == 0 {
			s = "Start!"
		}
		if i < 6 {
			ircutil.SendResponse(client, message.Source, message.Target,
				s)
		}
		if i == 0 {
			break
		}
		i--
	}
	ticker.Stop()
}
