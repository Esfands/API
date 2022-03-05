# API
*Still a work in progress*

## Routes
### Base URL
*Not published yet*

## Commands
#### Fetch all commands
- **Method:**
`GET`
- **URL:**
`/commands` 
- **Params:**
None

#### Get command by name
- **Method:**
`GET`
- **URL:**
`/commands/:name` 
- **Params:**
None

## On the Fly Commands
#### Fetch all OTF commands
- **Method:**
`GET`
- **URL:**
`/otf`
- **Params:**
None

#### Fetch OTF command by name
- **Method:**
`GET`
- **URL:**
`/otf/:name`
- **Params:**
None

## Feedback
#### Get all feedback
- **Method:**
`GET`
- **URL:**
`/feedback`
- **Params:**
None

#### Get feedback by ID
- **Method:**
`GET`
- **URL:**
`/feedback/:id`
- **Params:**
None

## Subathon Stats
### Most active chatters
- **Method:**
`GET`
- **URL:**
`/subathon/chatters`
- **Params:**
None

### Gifted Subs
- **Method:**
`GET`
- **URL:**
`/subathon/giftedsubs`
- **Params:**
None

### Bits Donated
- **Method:**
`GET`
- **URL:**
`/subathon/bitsdonated`
- **Params:**
None

## Twitch
### Get a user's information
*Only supports string for now, will add support for ID if user isn't a string*
- **Method:** 
`GET` 
- **URL:**
`/twitch/id`
- **Params** 
    **Required:** `user=[username]`

### Get a user's channel sub/bits/followers only emotes
*Supports IDs only for now*
- **Method:** 
`GET` 
- **URL:**
`/twitch/emotes`
- **Params** 
    **Required:** `user=[1234567890]`

### EventSub
Subscribed to the following events
- **Channel Information**
  - [stream.online](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#streamonline)
  - [stream.offline](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#streamoffline)
  - [channel.update](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelupdate)
- **Predictions**
  - [channel.prediction.begin](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpredictionbegin)
  - [channel.prediction.progress](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpredictionprogress)
  - [channel.prediction.lock](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpredictionlock)
  - [channel.prediction.end](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpredictionend)
- **Polls**
  - [channel.poll.begin](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpollbegin)
  - [channel.poll.progress](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpollprogress)
  - [channel.poll.end](https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types#channelpollprogress)

### Ideas
- Make a bot that logs users in the most active chats
- Keep track of a users chat badges