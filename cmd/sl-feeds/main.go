package main

import (
	"fmt"
	"io/ioutil"
	"os"

	_ "../../fetch"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func main() {
	config := Config{}

	app := cli.NewApp()
	app.Name = "sl-feeds"
	app.Usage = "Transform slackware ChangeLog.txt into RSS feeds"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
		cli.StringFlag{
			Name:  "dest, d",
			Usage: "Output RSS files to `DIR`",
		},
		cli.BoolFlag{
			Name:  "sample-config",
			Usage: "Output sample config file to stdout",
		},
	}

	// This is the main/default application
	app.Action = func(c *cli.Context) error {
		if c.Bool("sample-config") {
			c := Config{
				Dest: "$HOME/public_html/feeds/",
				Mirrors: []Mirror{
					Mirror{
						URL: "http://slackware.osuosl.org/",
						Releases: []string{
							"slackware-14.0",
							"slackware-14.1",
							"slackware-14.2",
							"slackware-current",
							"slackware64-14.0",
							"slackware64-14.1",
							"slackware64-14.2",
							"slackware64-current",
						},
					},
					Mirror{
						URL: "http://ftp.arm.slackware.com/slackwarearm/",
						Releases: []string{
							"slackwarearm-14.1",
							"slackwarearm-14.2",
							"slackwarearm-current",
						},
					},
				},
			}
			toml.NewEncoder(os.Stdout).Encode(c)
			return nil
		}

		fmt.Println(os.ExpandEnv(config.Dest))
		/*
			for each mirror in Mirrors
				if there is not a $release.RSS file, then fetch the whole ChangeLog
				if there is a $release.RSS file, then stat the file and only fetch remote if it is newer than the local RSS file
				if the remote returns any error (404, 503, etc) then print a warning but continue
		*/
		return nil
	}

	app.Before = func(c *cli.Context) error {
		if c.String("config") == "" {
			return nil
		}

		data, err := ioutil.ReadFile(c.String("config"))
		if err != nil {
			return err
		}
		if _, err := toml.Decode(string(data), &config); err != nil {
			return err
		}
		if c.String("dest") != "" {
			config.Dest = c.String("dest")
		}
		return nil
	}

	app.Run(os.Args)
}

type Config struct {
	Dest    string
	Mirrors []Mirror
}

type Mirror struct {
	URL      string
	Releases []string
}
