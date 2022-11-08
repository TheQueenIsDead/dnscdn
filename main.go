package main

import (
	"dnscdn/commands"
	"dnscdn/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	// Set default log level
	log.SetLevel(config.DsnCdnDefaultLogLevel)

	var file string
	var domain string

	app := &cli.App{
		Name:  "DNSCDN",
		Usage: "Store and retrieve media by use of 'free' DNS storage.",
		Action: func(*cli.Context) error {
			log.Info("Default action")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "file",
				//Value:       "",
				Usage:       "File to retrieve or upload, including extension.",
				Aliases:     []string{"f"},
				Destination: &file,
			},
			&cli.StringFlag{
				Name:        "domain",
				Usage:       "Domain to retrieve from or upload to.",
				Aliases:     []string{"d"},
				Destination: &domain,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "upload",
				Usage:  "Upload a file to a given DNS provider by means of TXT record.",
				Action: commands.UploadCommand,
			},
			{
				Name:   "download",
				Usage:  "Retrieve file data from DNS and save it locally.",
				Action: commands.DownloadCommand,
			},
			{
				Name:    "list",
				Aliases: []string{"d"},
				Usage:   "Given a domain, enumerate for previously saved DNSCDN media.",
				Action:  commands.ListCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
