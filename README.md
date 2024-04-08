# Shards [![Build Status](https://travis-ci.com/servusDei2018/shards.svg?branch=master)](https://app.travis-ci.com/github/servusdei2018/shards) [![CodeFactor](https://www.codefactor.io/repository/github/servusdei2018/shards/badge)](https://www.codefactor.io/repository/github/servusdei2018/shards) [![Go Reference](https://pkg.go.dev/badge/github.com/servusdei2018/shards.svg)](https://pkg.go.dev/github.com/servusdei2018/shards/v2)

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

## Buy me a Coffee

<a href="https://www.buymeacoffee.com/nbracy" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-red.png" alt="Buy Me A Coffee" style="height: 40px !important;width: 144px !important;" ></a>

## License
```
The MIT License (MIT)

Copyright (c) 2023-present The Shards Authors

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
