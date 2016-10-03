package utilcmd

import (
	"strings"

	"github.com/jasonpuglisi/ircutil"
)

// Say sends a message to a target. Function key: inami/utilcmd.Say
func Say(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	ircutil.SendPrivmsg(client, message.Args[0], strings.Join(message.Args[1:],
		" "))
}

// Init adds utilcmd's functions to the command map.
func Init(cmdMap ircutil.CmdMap) {
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Say", Say)
}
