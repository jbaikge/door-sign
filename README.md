# Automatic Door Sign

**Problem:** your family has no idea when you are on a meeting and often opens
the door in the midst of a riveting TPS report discussion.

**Solution:** Monitor powerd on MacOS for the screensaver disabling log entries
and contact Home Assistant to activate or deactivate a Z-Wave plug (seen as a
switch).

## Usage

```bash
$ go install github.com/jbaikge/door-sign
$ door-sign -api-url https://my.home-assistant.com/api -entity-id switch.door_sign_plug -token the-api-token-from-home-assistant
```

## Caveats

1. _video-playing_ does not activate the switch. This actually triggers on YouTube when hovering over thumbnails
2. Chime calls do not trigger any _powerd_ events
3. Teams calls work great
4. Slack calls work great
5. Zoom calls in the browser do not trigger any _powerd_ events
