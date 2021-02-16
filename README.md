# sl-feeds

This is for proccessing Slackware ChangeLog.txt -> RSS feeds that folks can
subscribed to.

Ultimately ending up at http://www.slackware.com/~vbatts/feeds/ or http://mirrors.slackware.com/feeds/

## Usage

```bash
go get github.com/vbatts/sl-feeds/cmd/sl-feeds
```

Create a configuration from the sample

```bash
sl-feeds --sample-config > ~/.sl-feeds.toml
```

crontab like:

```
0 */2 * * * ~/bin/sl-feeds -c ~/.sl-feeds.toml -q || mail -s "[sl-feeds] failed $(date +%D)" me@example.com
```
