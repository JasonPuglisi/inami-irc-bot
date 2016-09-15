package main

import (
  "fmt"

  "github.com/jasonpuglisi/ircutil"
)

// main creates a server and user, and establishes a connected. It starts a
// loop that runs until the client is stopped.
func main() {
  // Create a server for a connection
  server, err := ircutil.CreateServer("irc.rizon.net", 6667, true, "");
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Create a user for a connection
  user, err := ircutil.CreateUser("Inami", "inami", "Mahiru Inami", 8)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Establish a connection using a user and server
  client, err := ircutil.EstablishConnection(server, user);
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Loop until client is stopped
  for client.Active {}
}
