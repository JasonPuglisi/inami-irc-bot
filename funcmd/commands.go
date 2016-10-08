package funcmd

import (
	"math/rand"

	"github.com/jasonpuglisi/ircutil"
)

// Init adds funcmd's functions to the command map.
func Init(cmdMap ircutil.CmdMap) {
	ircutil.AddCommand(cmdMap, "inami/funcmd.EightBall", EightBall)
}

// EightBall sends a random magic 8-ball response.
// Function key: inami/funcmd.EightBall
func EightBall(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	responses := []string{"It is certain", "It is decidedly so",
		"Without a doubt", "Yes, definitely", "You may rely on it",
		"As I see it, yes", "Most likely", "Outlook good", "Yes",
		"Signs point to yes", "Reply hazy try again", "Ask again later",
		"Better not tell you now", "Cannot predict now",
		"Concentrate and ask again", "Don't count on it", "My reply is no",
		"My sources say no", "Outlook not so good", "Very doubtful"}
	ircutil.SendResponse(client, message.Source, message.Target,
		responses[rand.Intn(len(responses))])
}
