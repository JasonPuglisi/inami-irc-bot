package utilcmd

import (
	"fmt"
	"strings"

	"github.com/jasonpuglisi/inami-irc-bot/configutil"
	"github.com/jasonpuglisi/ircutil"
)

// SetProfileItem saves a user's profile item to persistent data.
func SetProfileItem(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Set name and value, and update scope/owner to match command scope.
	name, value := message.Args[0], strings.Join(message.Args[1:], " ")
	keys := []string{"", "", "utility/profile", name}
	configutil.UpdateScope(keys, message.Source, message.Source)

	// Set profile item in persistent data and send response with confirmation.
	configutil.SetValue(client, keys, value)
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("Your %s is now %s", name, value))
}

// GetProfileItem outputs a user's profile item saved to persistent data.
func GetProfileItem(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	// Set name and update scope/owner to match command scope.
	name := message.Args[0]
	keys := []string{"", "", "utility/profile", name}
	configutil.UpdateScope(keys, message.Source, message.Source)

	// Get profile item from name in persistent data.
	value, err := configutil.GetValue(client, keys)
	if err != nil {
		ircutil.Log(client, err.Error())
		ircutil.SendResponse(client, message.Source, message.Target,
			"Error getting profile item, try again later")
		return
	}
	if len(value) < 1 {
		ircutil.SendResponse(client, message.Source, message.Target,
			"Profile item not found, make sure it exists")
		return
	}

	// Send response with profile item.
	ircutil.SendResponse(client, message.Source, message.Target,
		fmt.Sprintf("Your %s is %s", name, value))
}
