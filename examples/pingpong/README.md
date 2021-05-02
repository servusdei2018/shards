## DiscordGo/Shards Ping Pong Example

This example demonstrates how to utilize [DiscordGo](https://github.com/bwmarrin/discordgo)
and [Shards](https://github.com/servusDei2018/shards) to create an extremely
scalable Ping Pong Bot.

This Bot will respond to "ping" with "Pong!" and "pong" with "Ping!".
This Bot will also respond to "restart" by performing a zero-downtime
rescaling restart. Simply enter this command to see it restart live,
without going offline.

**Open an issue on [Shards](https://github.com/servusDei2018/shards) if you are
having difficulties, or, join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.


From within the pingpong example folder, run the below command to
compile the example.

```sh
go build
```

### Usage

This example uses bot tokens for authentication only.
While user/password is supported by DiscordGo, it is not recommended for
bots.

```console
$ ./pingpong --help
Usage of ./pingpong:
  -t string
        Bot Token
```

The below example shows how to start the bot:

```console
$ ./pingpong -t YOUR_BOT_TOKEN
[INFO] Starting shard manager...
[INFO] Shard #0 connected.
[SUCCESS] Bot is now running.  Press CTRL-C to exit.
```
