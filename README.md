# ChatMeterBot 

[![Build Status](https://travis-ci.org/fjorgemota/chatmeterbot.svg?branch=master)](https://travis-ci.org/fjorgemota/chatmeterbot)

This is a bot, created for usage with Telegram, that automatically sends a sticker 'SILENCE' after 10 minutes of inactivity, and a sticker 'RIP CHAT' after 30 minutes of activity.

Note that these intervals are not precise, so it may not correspond exactly to 10 and 30 minutes, respectively, because the polling occurs in exactly these intervals. So it is possible to a sticker to be sent between 10 and 20 minutes in the case of the 'SILENCE' sticker, and between 30 and 60 minutes in the case of the 'RIP CHAT' sticker.

An important observation is that, additional to the rule of inactivity in a group, there is another rule that makes the bot to only send a sticker after a moment of intense activity in the group, with more than 10 messages sent for each minute in the interval counted. So, if an group sends more than 20 messages in 10 minutes, the 'SILENCE' sticker is sent after at least 10 minutes of inactivity, and so on.

You may add it to a group using it's username, [@ChatMeterBot](https://telegram.me/ChatMeterBot), note that it only works in groups and supergroups, not in the private mode (because it does not make sense to..).

## Installation

You may install it using **Go**:

```
$ go get -u github.com/fjorgemota/chatmeterbot
$ BOT_TOKEN=<your bot token informed by BotFather> chatmeterbot
```

Or using [Docker](http://docker.com):

```
docker run -e BOT_TOKEN=<your bot token informed by BotFather> fjorgemota/chatmeterbot
```
