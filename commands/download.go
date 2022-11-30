package commands

import (
	"dnscdn/lib"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"os"
)

// DownloadCommand checks to see if the given file is in the index, then enumerates the data TXT records based on
// the count stored in the index at the time.
// It then writes the file out to the current directory with the prefix "DNS-".
func DownloadCommand(cCtx *cli.Context) error {

	fileName := cCtx.String("file")
	domainName := cCtx.String("domain")
	idxFqdn := lib.IndexFqdn(domainName)

	// Lookup index record in order to find count of data records
	idx, err := net.LookupTXT(idxFqdn)
	if err != nil {
		return err
	}
	records := lib.ParseIndex(idx[0])

	// Check the index contains the file we're trying to retrieve
	if _, ok := records[fileName]; !ok {
		return errors.New(fmt.Sprintf("Could not locate file '%s' in domain index", fileName))
	}
	idxN := records[fileName]

	// Enumerate and write
	log.Infof("Retrieving %s (%d records)", fileName, idxN)

	fileData, err := dnsToFile(fileName, domainName, idxN)
	if err != nil {
		return err
	}

	log.Info("Data retrieved.")
	err = os.WriteFile(fmt.Sprintf("DNS-%s", fileName), fileData, os.ModeAppend)

	return err

}

func dnsToFile(fileName string, domainName string, idxN int) ([]byte, error) {

	// Enumerate records and store data in variable
	sfile := ""
	for i := 0; i < idxN; i++ {
		lookup := lib.DataFqdn(fileName, domainName, i)
		txt, err := net.LookupTXT(lookup)
		if err != nil {
			return nil, err
		}
		sfile = sfile + txt[0]
	}

	file, err := base64.StdEncoding.DecodeString(sfile)

	return file, err

}
