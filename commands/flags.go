package commands

import "github.com/urfave/cli/v2"

func DomainNameFlag(domain *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "domain",
		Usage:       "Domain to retrieve from or upload to.",
		Aliases:     []string{"d"},
		Destination: domain,
	}
}

func FileNameFlag(file *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "file",
		Usage:       "File to retrieve or upload, including extension.",
		Aliases:     []string{"f"},
		Destination: file,
	}
}
