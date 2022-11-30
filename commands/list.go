package commands

import (
	"dnscdn/config"
	"dnscdn/lib"
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
	"strconv"
	"strings"
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
	records := strings.Split(idx[0], config.DnsCdnIndexSeparator)

	// Loop records and print in table
	tblTemplate := "%-20s %-10s %-10s\n"
	fmt.Printf("DNSCDN files on %s:\n", domainName)
	fmt.Printf(tblTemplate, "Filename", "Length", "Size (kB)")
	for _, r := range records {
		recordComponents := strings.Split(r, ",") // filename,length

		filename := recordComponents[0]
		lengthStr := recordComponents[1]
		lengthInt, _ := strconv.Atoi(lengthStr)
		sizeKb := (lengthInt * config.DnsCdnDataSize) / 1000

		fmt.Printf(tblTemplate, filename, lengthStr, fmt.Sprintf("%d", sizeKb))
	}

	return err
}
