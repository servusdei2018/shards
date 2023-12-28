# Shards [![Build Status](https://travis-ci.com/servusDei2018/shards.svg?branch=master)](https://travis-ci.com/servusDei2018/shards) [![CodeFactor](https://www.codefactor.io/repository/github/servusdei2018/shards/badge)](https://www.codefactor.io/repository/github/servusdei2018/shards) [![Go Reference](https://pkg.go.dev/badge/github.com/servusdei2018/shards.svg)](https://pkg.go.dev/github.com/servusdei2018/shards/v2)

Configurable, scalable and automatic sharding library for `discordgo`.

## Features
 - Automatic scaling: head-ache free configuration runs out of the box.
 - [Zero-downtime restarts](https://pkg.go.dev/github.com/servusdei2018/shards#Manager.Restart): make downtime a thing of the past.
 - [Slash commands](https://pkg.go.dev/github.com/servusdei2018/shards#Manager.ApplicationCommandCreate): integrate the latest in Discord functionality.

### Installing

This assumes you already have a working Go environment, if not please [install Go](https://golang.org/doc/install) first.

```sh
go get github.com/servusdei2018/shards/v2
```

### Usage

Import the package into your project, like so:

```go
import "github.com/servusdei2018/shards/v2"
```