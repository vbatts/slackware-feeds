# sl-feeds

This is for proccessing Slackware ChangeLog.txt -> RSS feeds that folks can
subscribed to.

Ultimately ending up at http://www.slackware.com/~vbatts/feeds/

## Usage

```bash
go get github.com/vbatts/sl-feeds
```

crontab like:

```
0 */2 * * * ~/bin/sl-feeds -q || echo "$(date): failed to poll changelogs" | mail -s "[slackrss] changelog_http_poll failed $(date +%D)" me@example.com
```
