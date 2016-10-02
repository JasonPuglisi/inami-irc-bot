package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jasonpuglisi/ircutil"
)

// main requests a server and user which it uses establish a connection. It
// runs a loop to keep the client alive until it is no longer active.
func main() {
	// Set config and debug flags, then parse command line arguments.
	configPtr := flag.String("config", "config.json", "configuration file")
	debugPtr := flag.Bool("debug", false, "debugging mode")
	flag.Parse()

	// Get configuration from filename.
	config, err := getConfig(*configPtr)
	if err != nil {
		fmt.Printf("Error opening %s, make sure the file exists.\n%s\n",
			*configPtr, err)
		return
	}

	// Seed random number generator.
	rand.Seed(time.Now().UnixNano())

	// Declare slice to store clients.
	var clients []*ircutil.Client

	// Loop through all clients in config to establish their connections.
	for i := range config.Clients {
		client := &config.Clients[i]

		// Get server from config and reference it in client.
		server, err := getServer(config, client.ServerID)
		if err != nil {
			fmt.Printf("Error getting server %s, make sure it exists in %s.\n%s\n",
				client.ServerID, *configPtr, err)
			return
		}
		client.Server = server

		// Get user from config and reference it in client.
		user, err := getUser(config, client.UserID)
		if err != nil {
			fmt.Printf("Error getting user %s, make sure it exists in %s.\n%s\n",
				client.UserID, *configPtr, err)
			return
		}
		client.User = user

		// Set debugging mode, ready function, and done channel for client.
		client.Debug = *debugPtr
		client.Ready = Init
		client.Done = make(chan bool, 1)
		client.Nick = client.User.Nick

		// Establish a connection with the created client.
		err = ircutil.EstablishConnection(client)
		if err != nil {
			fmt.Printf("Error establishing connection with %s/%s, make sure its "+
				"settings are valid in %s.\n%s\n", client.ServerID, client.UserID,
				*configPtr, err)
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
