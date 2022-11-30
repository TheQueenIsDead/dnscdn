package commands

import (
	"dnscdn/lib"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"os"
	"strconv"
	"strings"
)

func DownloadCommand(cCtx *cli.Context) error {

	fileName := cCtx.String("file")
	domainName := cCtx.String("domain")
	idxFqdn := lib.IndexFqdn(domainName)

	log.Infof("Retrieving %s.", fileName)

	fileData, err := dnsToFile(idxFqdn)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("DNS-%s", fileName), fileData, os.ModeAppend)

	return err

}

func dnsToFile(idxFqdn string) ([]byte, error) {

	mediaSplit := strings.Split(idxFqdn, ".media.")
	fileName := mediaSplit[0]
	domainName := mediaSplit[len(mediaSplit)-1]

	logger := log.StandardLogger().WithFields(log.Fields{
		"fqdn":   idxFqdn,
		"file":   fileName,
		"domain": domainName,
	})

	// Lookup index record in order to find count of records
	idx, err := net.LookupTXT(idxFqdn)
	if err != nil {
		return nil, err
	}
	idxN, _ := strconv.Atoi(idx[0])

	// Enumerate records and store data in variable
	sfile := ""
	for i := 0; i < idxN; i++ {
		lookup := fmt.Sprintf("%s.%d.media.%s", fileName, i, domainName)
		logger = logger.WithField("record", lookup)
		logger.Debugf("Retrieving TXT record.")
		txt, err := net.LookupTXT(lookup)
		if err != nil {
			logger.WithError(err).Error("Could not retrieve TXT record.")
			return nil, err
		}
		sfile = sfile + txt[0]
	}

	file, err := base64.StdEncoding.DecodeString(sfile)

	return file, err

}
