package main

import (
  "flag"
  "fmt"

  "github.com/jasonpuglisi/ircutil"
)

// TODO: All the static values in here will be replaced with config-loaded
//       values. This includes server/user info and commands.

// main requests a server and user which it uses establish a connection. It
// runs a loop to keep the client alive until it is no longer active.
func main() {
  // Set config and debug flags, then parse command line arguments.
  //configPtr := flag.String("config", "config.json", "configuration file")
  debugPtr := flag.Bool("debug", false, "debugging mode")
  flag.Parse()

  // Request a server with the specified details.
  server, err := ircutil.CreateServer("irc.rizon.net", 6697, true, "")
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Request a user with the specified details.
  user, err := ircutil.CreateUser("Inami", "inami", "Mahiru Inami", 8)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Establish a connection and get a client using user and server details as
  // well as an initialization function and debugging setting.
  client, err := ircutil.EstablishConnection(server, user, Init, *debugPtr)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Loop until client is no longer active.
  for client.Active {}
}

// Init is executed after the client it connected and registered to the server.
func Init(client *ircutil.Client) {
  //ircutil.Join(client, "#snowie", "")
}
