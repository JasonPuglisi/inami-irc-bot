package animecmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jasonpuglisi/inami-irc-bot/configutil"
	"github.com/jasonpuglisi/ircutil"
)

// Init adds animecmd's functions to the command map.
func Init(cmdMap ircutil.CmdMap) {
	ircutil.AddCommand(cmdMap, "inami/animecmd.Countdown", Countdown)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Alias", Alias)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Search", Search)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Watch", Watch)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Progress", Progress)
	ircutil.AddCommand(cmdMap, "inami/animecmd.Next", Next)
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

// Alias saves a show identifier as a custom alias.
func Alias(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Set alias and value, and update scope/owner to match command scope.
	id, alias := message.Args[0], message.Args[1]
	keys := []string{"", "", "anime/shows", alias}
	configutil.UpdateScope(keys, message.Source, message.Target)

	// Set alias, intialize episode progress, and send response with
	// confirmation.
	configutil.SetValue(client, keys, id)
	keys[2] = "anime/progress"
	configutil.SetValue(client, keys, "0")
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("Aliased %s to %s", id, alias))
}

// Search searches an anime database and returns relevant search results.
// Function key: inami/animecmd.Search
func Search(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Search Kitsu for shows matching query.
	shows, err := search(strings.Join(message.Args, " "))
	if err != nil {
		ircutil.Log(client, err.Error())
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
			fmt.Sprintf("- [%s] %s", s.ID, s.Attributes.Title))
	}
}

// Watch gets data of a show and starts a group watching session.
// Function key: inami/animecmd.Watch
func Watch(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Set alias and update scope/owner to match command scope.
	alias := message.Args[0]
	keys := []string{"", "", "anime/shows", alias}
	configutil.UpdateScope(keys, message.Source, message.Target)

	// Get show id from alias in persistent data.
	id, err := configutil.GetValue(client, keys)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error getting alias, try again later")
		return
	}
	if len(id) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			"Alias not found, make sure you've assigned a show to it")
		return
	}

	// Get episode number from alias in persistent data.
	keys[2] = "anime/progress"
	numStr, _ := configutil.GetValue(client, keys)
	num, _ := strconv.Atoi(numStr)
	num++

	// Get show data from Kitsu.
	show, err := show(id)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error fetching show, try again later")
		return
	}

	// Send response if show not found.
	if len(show.Attributes.Slug) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			fmt.Sprintf("Show not found, make sure your alias is using %s",
				"the name after /anime/ in the show's URL"))
		return
	}

	// Get episode data from Kitsu.
	episodes, err := episodes(id)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error fetching episodes, try again later")
		return
	}

	// Find specific episode.
	var episode *result
	for i := range episodes {
		e := &episodes[i]
		if e.Attributes.Number == num {
			episode = e
      break
		}
	}

	// Clear episode title if it doesn't exist.
	episodeTitle := ""
	if episode != nil && len(episode.Attributes.Title) > 0 {
		episodeTitle = fmt.Sprintf(" \"%s\"", episode.Attributes.Title)
	}

	// Start countdown and sleep before sending episode title.
	Countdown(client, command, message)
	time.Sleep(time.Second * 5)

	// Send response with episode information, and increment episode number.
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("You're watching %s Episode %d%s", show.Attributes.Title, num,
			episodeTitle))
	configutil.SetValue(client, keys, strconv.Itoa(num))
}

// Progress updates a show's episode number.
func Progress(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Parse episode number from arguments.
	num, err := strconv.Atoi(message.Args[1])
	if err != nil || num < 0 {
		ircutil.SendResponse(client, message.Source, message.Target,
			"Invalid episode number")
	}

	// Set alias and update scope/owner to match command scope.
	alias := message.Args[0]
	keys := []string{"", "", "anime/shows", alias}
	configutil.UpdateScope(keys, message.Source, message.Target)

	// Get show name from alias in persistent data to make sure it exists.
	id, err := configutil.GetValue(client, keys)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error checking for alias, try again later")
		return
	}
	if len(id) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			"Alias not found, make sure you've assigned a show to it")
		return
	}

	// Set episode number in persistent data and send response with confirmation.
	keys[2] = "anime/progress"
	plural := "s"
	if num == 1 {
		plural = ""
	}
	configutil.SetValue(client, keys, strconv.Itoa(num))
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("Updated show progress, you've watched %d episode%s", num,
			plural))
}

// Next returns information about the next episode for a show.
func Next(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Set alias and update scope/owner to match command scope.
	alias := message.Args[0]
	keys := []string{"", "", "anime/shows", alias}
	configutil.UpdateScope(keys, message.Source, message.Target)

	// Get show id from alias in persistent data to make sure it exists.
	id, err := configutil.GetValue(client, keys)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error checking for alias, try again later")
		return
	}
	if len(id) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			"Alias not found, make sure you've assigned a show to it")
		return
	}

	// Get episode number from alias in persistent data.
	keys[2] = "anime/progress"
	strNum, _ := configutil.GetValue(client, keys)
	num, _ := strconv.Atoi(strNum)
	num++

	// Get show data from Kitsu.
	show, err := show(id)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error fetching show, try again later")
		return
	}

	// Send response if show not found.
	if len(show.Attributes.Slug) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			fmt.Sprintf("Show not found, make sure your alias is using %s",
				"the name after /anime/ in the show's URL"))
		return
	}

	// Get episode data from Kitsu.
	episodes, err := episodes(id)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error fetching episodes, try again later")
		return
	}

	// Find specific episode.
	var episode *result
	for i := range episodes {
		e := &episodes[i]
		if e.Attributes.Number == num {
			episode = e
      break
		}
	}

	// Clear episode title if it doesn't exist.
	episodeTitle := ""
	if episode != nil && len(episode.Attributes.Title) > 0 {
		episodeTitle = fmt.Sprintf(" \"%s\"", episode.Attributes.Title)
	}

	// Send response with next episode information.
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("Next up for %s is Episode %d%s", show.Attributes.Title, num,
			episodeTitle))
}
