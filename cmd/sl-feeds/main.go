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

	app.Action = func(c *cli.Context) error {
		if c.Bool("sample-config") {
			c := Config{
				Dest: "./public_html/feeds/",
			}
			toml.NewEncoder(os.Stdout).Encode(c)
			return nil
		}

		fmt.Println("boom! I say!")
		fmt.Println(config)
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
		return nil
	}

	app.Run(os.Args)
}

type Config struct {
	Dest string
}
