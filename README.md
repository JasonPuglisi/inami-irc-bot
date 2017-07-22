# Inami IRC Bot

Modular multi-network IRC bot with custom commands and persistent user data.

Developed along with [ircutil](https://github.com/JasonPuglisi/ircutil). Check
out the readme for more information.

## Dependencies

- [ircutil](https://github.com/JasonPuglisi/ircutil)

## Optional Dependencies

- [gomemcache](https://github.com/bradfitz/gomemcache)

## Usage

Run `go get github.com/jasonpuglisi/inami-irc-bot` to grab the program. Then
run `go install` in the root directory to install the executable.

You must create the configuration file `config.json` before the bot will work.
An example configuration file is provided as `config-example.json`. This file
contains all possible commands, and custom commands can be added in the same
format. The `servers` and `users` sections both contain `id` fields that should
be referenced in the `clients` section (these are your IRC server connections).
Most options, such as passwords, can be omitted or left blank (they are in the
example for reference).

Some functions may benefit from having `gomemcache` installed, but they will
work without it. To use this dependency, you must have
[`memcached`](https://memcached.org/) installed. Installing this dependency is
recommended to reduce API calls and improve command response time.

## Overview

Supports connections to an unlimited number of IRC servers, and interfaces
with the dependency `ircutil` to distinguish log messages from different
clients. Is able to save persistent user data in various scopes that can be
accessed later.

Focuses on extensibility. Note that because of the nature of Go programs, you
must modify an existing source file to add your own functionality. See the
section below for details.

Comes with a number of administrative commands that can make a bot change its
nickname, speak in a channel, and more. Includes a module for _fun_ commands,
which may be extended in the future.

The main feature is currently a system for searching an anime database for
shows, and storing show progress for a channel. This is meant to coordinate
group watching of shows within a channel. Users can query what the next episode
is, or start a countdown to synchronize watching. It supports fetching and
displaying individual episode titles if they exist with the database. This
feature currently uses the [Kitsu Edge API](http://docs.kitsu17.apiary.io/).

## Extending

This bot is built with extensibility in mind, and you can add modules for your
own use by following the same format as the existing modules (such as those in
[`animecmd`](animecmd) and [`funcmd`](funcmd)). Note that **in order for new
modules to run**, you must add them to the `import` statement at the top of
[`client.go`](client.go). New modules cannot be dynamically loaded from a
folder due to the nature of Go, so they must be imported statically.

Keep in mind that [`client.go`](client.go) is checked into the source
repository. You may need to discard your changes before pulling an updated
version of the file, and restore them after. If you believe your module would
benefit other users of the project, you are encouraged to submit a pull
request, in which case your module import would remain in
[`client.go`](client.go).
