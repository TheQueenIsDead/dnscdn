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
		Commands: []*cli.Command{
			{
				Name:   "upload",
				Usage:  "Upload a file to a given DNS provider by means of TXT record.",
				Action: commands.UploadCommand,
				Flags: []cli.Flag{
					commands.FileNameFlag(&file),
					commands.DomainNameFlag(&domain),
				},
			},
			{
				Name:   "download",
				Usage:  "Retrieve file data from DNS and save it locally.",
				Action: commands.DownloadCommand,
				Flags: []cli.Flag{
					commands.FileNameFlag(&file),
					commands.DomainNameFlag(&domain),
				},
			},
			{
				Name:   "delete",
				Usage:  "Remove all file specific records for a domain.",
				Action: commands.DeleteCommand,
				Flags: []cli.Flag{
					commands.FileNameFlag(&file),
					commands.DomainNameFlag(&domain),
				},
			},
			{
				Name:   "list",
				Usage:  "Given a domain, enumerate for previously saved DNSCDN media.",
				Action: commands.ListCommand,
				Flags: []cli.Flag{
					commands.DomainNameFlag(&domain),
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
