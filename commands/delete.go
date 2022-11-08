package commands

import (
	"dnscdn/provider"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"strconv"
)

func DeleteCommand(cCtx *cli.Context) error {

	fileName := cCtx.String("file")
	domainName := cCtx.String("domain")

	//mediaSplit := strings.Split(idxFqdn, ".media.")
	//fileName := mediaSplit[0]
	//domainName := mediaSplit[len(mediaSplit)-1]
	//

	dnsProvider := provider.CloudflareDnsProvider{}
	err := dnsProvider.Authenticate()
	if err != nil {
		log.WithError(err).Error("Error encountered calling Authenticate()")
		return nil
	}

	idxFqdn := fmt.Sprintf("%s.media.%s", fileName, domainName)

	// Lookup index record in order to find count of records
	idx, err := net.LookupTXT(idxFqdn)
	if err != nil {
		return err
	}
	idxN, _ := strconv.Atoi(idx[0])

	// Enumerate records and store data in variable
	for i := 0; i < idxN; i++ {
		fqdn := fmt.Sprintf("%s.%d.media.%s", fileName, i, domainName)
		err := dnsProvider.DeleteRecord(fqdn)
		if err != nil {
			return err
		}
	}

	// Cleanup index
	err = dnsProvider.DeleteRecord(idxFqdn)
	if err != nil {
		return err
	}

	return err
}
