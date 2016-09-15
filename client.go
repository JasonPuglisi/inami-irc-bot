package main

import (
  "fmt"

  "github.com/jasonpuglisi/ircutil"
)

func main() {
  server, err := ircutil.CreateServer("irc.rizon.net", 6667, true, "")
  if err != nil {
    fmt.Println(err)
    return
  }

  user, err := ircutil.CreateUser("Inami", "inami", "Mahiru Inami", 8)
  if err != nil {
    fmt.Println(err)
    return
  }

  // Create connection to IRC server
  ircutil.EstablishConnection(server, user);
}
