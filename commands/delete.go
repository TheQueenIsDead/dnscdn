package commands

import (
	"dnscdn/lib"
	"dnscdn/provider"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
)

// DeleteCommand iterates over TXT records containing file data and removes them, updating the DNSCDN index afterwards.
func DeleteCommand(cCtx *cli.Context) error {

	// Parse flags
	fileName := cCtx.String("file")
	domainName := cCtx.String("domain")

	// Lookup index record in order to find count of data records
	idxFqdn := lib.IndexFqdn(domainName)
	idx, err := net.LookupTXT(idxFqdn)
	if err != nil {
		log.Fatal("Could not retrieve DNSCDN TXT index record while deleting file. Perhaps no files have been uploaded yet?")
	}
	records := lib.ParseIndex(idx[0])

	// Authenticate
	dnsProvider := provider.CloudflareDnsProvider{}
	err = dnsProvider.Authenticate()
	if err != nil {
		log.WithError(err).Error("Error encountered calling Authenticate()")
		return nil
	}

	// Verify the index contains the file we're trying to delete
	if _, ok := records[fileName]; !ok {
		return errors.New(fmt.Sprintf("Could not locate file '%s' in domain index", fileName))
	}

	// Enumerate file records and delete
	idxN := records[fileName]
	for i := 0; i < idxN; i++ {
		fqdn := lib.DataFqdn(fileName, domainName, i)
		err := dnsProvider.DeleteRecord(fqdn)
		if err != nil {
			return err
		}
	}

	// Remove index entry and persist
	delete(records, fileName)
	err = lib.UpdateIndex(domainName, records, &dnsProvider)
	if err != nil {
		return err
	}

	return err
}
