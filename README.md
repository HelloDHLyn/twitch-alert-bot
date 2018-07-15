# Twitch Alert Bot

## Development

### Prerequisite

  - Go 1.X

### Environment Variables

  - `DISCORD_BOT_TOKEN`
  - `TWITCH_CLIENT_ID`
  - `LYNLAB_API_KEY`

### Run

```bash
go run main.go
```

### Deploy

```bash
# Build docker image
docker build -t twitch-alert-bot .

# Run
docker run -e '...' twitch-alert-bot
```
