package commands

import (
	"dnscdn/config"
	"dnscdn/lib"
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
	"strconv"
)

// ListCommand retrieves an index TXT record and parses it, displaying it to the user in order to list DNSCDN
// files stored on a domain.
func ListCommand(cCtx *cli.Context) error {

	// Parse flags
	domainName := cCtx.String("domain")
	indexFqdn := lib.IndexFqdn(domainName)

	// Retrieve index TXT record
	idx, err := net.LookupTXT(indexFqdn)
	if err != nil {
		return err
	}

	// Enumerate file records split on ";"
	records := lib.ParseIndex(idx[0])

	// Loop records and print in table
	tblTemplate := "%-20s %-10s %-10s\n"
	fmt.Printf("DNSCDN files on %s:\n", domainName)
	fmt.Printf(tblTemplate, "Filename", "Length", "Size (kB)")
	for filename, length := range records {
		fmt.Printf(tblTemplate, filename, strconv.Itoa(length), fmt.Sprintf("%d", length*config.DnsCdnDataSize/1024))
	}

	return err
}
