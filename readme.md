# Cord Bridge

Some ready to go functions to control a discord bot to _kinda_ act like a database

[npm package](https://www.npmjs.com/package/cordbridge)

# Functions

- **Create Categories**
- **Create Channels**
- **Find Channels**
- **Send Messages**
- **Read All Messages**
- **Read Last X Messages**
- **Edit Messages**
- **Delete Messages**

# Usage

## 🔹 Initializing the Bot

To start the bot, create a new client with your token and guild ID:

```go
package main

import (
    "log"
    "os"
    "github.com/yourusername/cordbridge"
)

func main() {
    token := os.Getenv("bot-token")
    guildID := "your-guild-id"

    client, err := cordbridge.NewClient(token, guildID)
    if err != nil {
        log.Fatal("Error creating client:", err)
    }

    err = client.Open()
    if err != nil {
        log.Fatal("Error opening connection:", err)
    }
    defer client.Close()
}
```

or use your own discordgo session

```go
session, _ := discordgo.New("Bot " + token)

client := &cordbridge.Client{
    Session: session,
    GuildID: "your-guild-id",
}
```

## 🔹 Creating a Category

```go
category, err := client.CreateCategory(categoryName)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Category created with ID: %s\n", category.ID)
```

## 🔹 Creating a Channel

```go
channel, err := client.CreateChannel(channelName, categoryID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Channel created with ID: %s\n", channel.ID)
```

## 🔹 Finding a Channel

```go
channel, err := client.FindChannelByName(channelName, categoryID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found channel ID: %s\n", channel.ID)
```

## 🔹 Sending Messages

```go
err := client.SendMessageToChannel(channelID, "Hello, World!")
if err != nil {
    log.Fatal(err)
}
```

## 🔹 Reading All Messages From a Channel

```go
allMessages, err := client.ReadAllMessagesFromChannel(channelID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Fetched %d messages\n", len(allMessages))
```

## 🔹 Reading Last X Messages From a Channel

```go
last10, err := client.ReadLastXMessagesFromChannel(channelID, 10)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Fetched last %d messages\n", len(last10))
```

## 🔹 Editing Messages

```go
err := client.EditMessageByID(channelID, messageID, newContent)
if err != nil {
    log.Fatal(err)
}
```

## 🔹 Deleting Messages

```go
err := client.DeleteMessageByID(channelID, messageID)
if err != nil {
    log.Fatal(err)
}
```

make sure to give the bot permissions to send and read msgs

why would someone use this? idk u can use it to store msgs for a chat app or smth but u have to keep rate limits in mind
