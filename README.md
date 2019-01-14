# dashboard-nerf

dashboard-nerf provides a web-dashboard to play MP3/WAV/MP4/WEBM sounds and movies,
either on the server/dashboard (to annoy your colleagues) or locally in your browser (for preview).

## installing

Grab and just unzip the latest [release](../../releases) for your OS.

## building from source

dashboard-nerf relies on golang (>= 1.7) standard library only. Build using `go get`:

```bash
go get github.com/schnoddelbotz/dashboard-nerf
```

Alternatively, clone this repository and run `make dependencies && make && make run`.
This will pull in `go-bindata`, which is used to update (bundled) `assets.go`
(i.e. the web application components). `make run` will also download example/test
media as mentioned in the next section.

## example invocation

Get some sounds and videos for playback purposes first! Below example assumes your
media files in `test_media` subfolder. Try with mine:
```bash
curl -s https://jan.hacker.ch/test_media.tgz | tar -xzf -
```

Obviously, you also need media players. On Debian/Ubuntu, try:
```bash
apt install sox libsox-fmt-mp3 vlc festival
```

Built-in defaults for players:

| Platform    | Videoplayer | Audioplayer | Speech    |
| ----------- | ----------- | ----------- | --------- |
| Darwin      | VLC.app     | afplay `*`  | say `*`   |
| Linux/amd64 | cvlc (VLC)  | play (sox)  | festival  |
| Linux/arm   | omxplayer   | play (sox)  | festival  |

Items marked with `*` are available by default (no installation required).
On RaspberryPi, `omxplayer` should be preferred over `vlc` (hardware acceleration);
unfortunately, it does not support playback of `webm` videos.

Finally, run dashboard-nerf:
```bash
  ./dashboard-nerf -media test_media
```

Overriding default players is possible by providing command line arguments:
```bash
  ./dashboard-nerf -media test_media \
    -audioplayer "mpg123 '%s'" \
    -videoplayer "mplayer '%s'" \
    -speech "espeak '%s'"
```
Ensure to include `'%s'` in your arguments -- it will be replaced by filename to be played / text to be spoken.

## hints

- For VLC, disable fullscreen controls once manually [via GUI](https://wiki.videolan.org/VSG:Usage:Controller/)
- For VLC, edit playlist settings (via "Show all" button in settings) to enable "Play and exit"

## license

MIT

(c) 2019 Jan Hacker and contributors
