package main

import (
  "flag"
  "fmt"
  "os"

  "github.com/jasonpuglisi/ircutil"
)

// TODO: All the static values in here will be replaced with config-loaded
//       values. This includes server/user info and commands.

// main requests a server and user which it uses establish a connection. It
// runs a loop to keep the client alive until it is no longer active.
func main() {
  // Set config and debug flags, then parse command line arguments.
  configPtr := flag.String("config", "config.json", "configuration file")
  debugPtr := flag.Bool("debug", false, "debugging mode")
  flag.Parse()

  // Attempt to open configuration file.
  _, err := os.Open(*configPtr)
  if err != nil {
    fmt.Printf("Unable to open configuration file \"%s\". Make sure the " +
      "file exists.\n[%s]\n", *configPtr, err)
    return
  }

  // Declare slice to store clients.
  var clients []*ircutil.Client

  // Request a server with the specified details.
  server, err := ircutil.CreateServer("irc.rizon.net", 6697, true, "")
  if err != nil {
    fmt.Println(err)
    return
  }

  // Request a user with the specified details.
  user, err := ircutil.CreateUser("Inami", "inami", "Mahiru Inami", "i")
  if err != nil {
    fmt.Println(err)
    return
  }

  // Establish a connection and get a client using user and server details as
  // well as an initialization function and debugging setting.
  client, err := ircutil.EstablishConnection(server, user, Init, *debugPtr)
  if err != nil {
    fmt.Println(err)
    return
  }

  // Add client to client slice.
  clients = append(clients, client)

  // Loop until all clients are no longer active.
  for _, c := range clients {
    <-c.Done
  }
}

// Init is executed after the client it connected and registered to the server.
func Init(client *ircutil.Client) {
  ircutil.SendJoin(client, "#inami", "")
}
