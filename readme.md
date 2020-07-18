# On Page Change Notification Bot

Easy to set-up discord bot or standalone app

# Build Instructions

```
go get github.com/bwmarrin/discordgo
go build
./OPCNB TIME_INTERVAL(int) DISCORD_BOT_TOKEN

TIME_INTERVAL = time in seconds between each check
Ex : 
./OPCNB 3 AB123kLM123mkfmi123jk
```

# Before running the bot

```
Create a file called "env.cfg" in which you specify the bot configuration.
Ex:
{"Prefix":"!","URL":"https://github.com/g0g05arui/OPCNB","Channel":"general","Online":true}
```

# Commanads
```
For commands I will use the default prefix (!)
!set URL
    updates the query URL
!start
    starts sending messages about the page
!stop
    stops sending messages about the page
!what
    shows the query URL
!chch abcd
    changes the channel where the messages are sent on to "abcd" (channel must be valid or bot won't work unless changed)
```