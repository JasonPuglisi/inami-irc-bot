package utilcmd

import (
	"fmt"
	"strings"

	"github.com/jasonpuglisi/ircutil"
)

// Init adds utilcmd's functions to the command map.
func Init(cmdMap ircutil.CmdMap) {
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Nick", Nick)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Join", Join)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Part", Part)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Say", Say)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Notify", Notify)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.Do", Do)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.GetProfileItem", GetProfileItem)
	ircutil.AddCommand(cmdMap, "inami/utilcmd.SetProfileItem", SetProfileItem)
}

// Nick updates a nickname. Function key: inami/utilcmd.Nick
func Nick(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	ircutil.SendNick(client, message.Args[0])
}

// Join attahces to a channel with an optional password.
// Function key: inami/utilcmd.Join
func Join(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	pass := ""
	if len(message.Args) > 1 {
		pass = message.Args[1]
	}
	ircutil.SendJoin(client, message.Args[0], pass)
}

// Part detaches from a channel. Function key: inami/utilcmd.Part
func Part(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	msg := ""
	if len(message.Args) > 1 {
		msg = strings.Join(message.Args[1:], " ")
	}
	ircutil.SendPart(client, message.Args[0], msg)
}

// Say sends a message to a target. Function key: inami/utilcmd.Say
func Say(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	ircutil.SendPrivmsg(client, message.Args[0], strings.Join(message.Args[1:],
		" "))
}

// Notify sends a notice to a target. Function key: inami/utilcmd.Notify
func Notify(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	ircutil.SendNotice(client, message.Args[0], strings.Join(message.Args[1:],
		" "))
}

// Do performs an action at a target. Function key: inami/utilcmd.Do
func Do(client *ircutil.Client, command *ircutil.Command,
	message *ircutil.Message) {
	ircutil.SendPrivmsg(client, message.Args[0], fmt.Sprintf("\x01ACTION %s\x01",
		strings.Join(message.Args[1:], " ")))
}
