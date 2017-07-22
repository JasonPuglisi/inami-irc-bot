package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jasonpuglisi/inami-irc-bot/animecmd"
	"github.com/jasonpuglisi/inami-irc-bot/configutil"
	"github.com/jasonpuglisi/inami-irc-bot/funcmd"
	"github.com/jasonpuglisi/inami-irc-bot/utilcmd"
	"github.com/jasonpuglisi/ircutil"
)

// main requests a server and user which it uses establish a connection. It
// runs a loop to keep the client alive until it is no longer active.
func main() {
	// Set config and debug flags, then parse command line arguments.
	configPtr := flag.String("config", "config.json", "configuration file")
	dataPtr := flag.String("data", "data.json", "data file")
	debugPtr := flag.Bool("debug", false, "debugging mode")
	flag.Parse()

	// Get configuration from filename.
	config, err := configutil.GetConfig(*configPtr)
	if err != nil {
		fmt.Printf("Error opening %s, %s.\n%s\n",
			*configPtr, "make sure the file exists and is correctly formatted", err)
		return
	}

	// Create data file if it doesn't already exist.
	_, err = os.Stat(*dataPtr)
	if os.IsNotExist(err) {
		err = ioutil.WriteFile(*dataPtr, []byte("{}"), 0644)
		if err != nil {
			fmt.Printf("Error creating %s.\n%s\n", *dataPtr, err)
			return
		}
	}

	// Get data from filename.
	data, err := configutil.GetData(*dataPtr)
	if err != nil {
		fmt.Printf("Error opening %s, %s.\n%s\n",
			*dataPtr, "make sure the file exists and is correctly formatted", err)
		return
	}

	// Seed random number generator.
	rand.Seed(time.Now().UnixNano())

	// Initialize and import commands.
	commands := config.Commands
	cmdMap := ircutil.InitCommands()
	utilcmd.Init(cmdMap)
	funcmd.Init(cmdMap)
	animecmd.Init(cmdMap)

	// Declare slice to store clients.
	var clients []*ircutil.Client

	// Loop through all clients in config to establish their connections.
	for i := range config.Clients {
		client := &config.Clients[i]

		// Get server from config and reference it in client.
		server, err := configutil.GetServer(config, client.ServerID)
		if err != nil {
			fmt.Printf("Error getting server %s, make sure it exists in %s.\n%s\n",
				client.ServerID, *configPtr, err)
			return
		}
		client.Server = server

		// Get user from config and reference it in client.
		user, err := configutil.GetUser(config, client.UserID)
		if err != nil {
			fmt.Printf("Error getting user %s, make sure it exists in %s.\n%s\n",
				client.UserID, *configPtr, err)
			return
		}
		client.User = user

		// Set client values.
		client.Data = data
		client.DataFile = dataPtr
		client.Commands = commands
		client.CmdMap = cmdMap
		client.Debug = *debugPtr
		client.Ready = Init
		client.Done = make(chan bool, 1)
		client.Nick = client.User.Nick

		// Establish a connection with the created client.
		err = ircutil.EstablishConnection(client)
		if err != nil {
			fmt.Printf("%s %s/%s, %s %s.\n%s\n",
				"Error establishing connection with", client.ServerID, client.UserID,
				"make sure its settings are valid in", *configPtr, err)
			return
		}

		// Add client to client slice.
		clients = append(clients, client)

		// Sleep for a second so we don't hit the same server too fast.
		if i < len(config.Clients)-1 {
			time.Sleep(time.Second)
		}
	}

	// Loop until all clients are no longer active.
	for _, c := range clients {
		<-c.Done
	}
}

// Init is executed after the client it connected and registered to the server.
func Init(client *ircutil.Client) {
	// Authenticate with Nickserv if a password is specified.
	if client.Nick == client.User.Nick && len(client.Authentication.Nickserv) >
		0 {
		ircutil.SendNickservPass(client, client.Authentication.Nickserv)
	}

	// Set user modes if specified.
	if len(client.Modes) > 0 {
		ircutil.SendModeUser(client, client.Modes)
	}

	// Join all of a client's channels.
	for i := range client.Channels {
		c := strings.Split(client.Channels[i], " ")
		pass := ""
		if len(c) > 1 {
			pass = c[1]
		}
		ircutil.SendJoin(client, c[0], pass)

		// Sleep for half a second so we don't join channels too fast.
		if i < len(client.Channels)-1 {
			time.Sleep(time.Millisecond * 500)
		}
	}
}
