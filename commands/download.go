package commands

import (
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"os"
	"strings"
)

func DownloadCommand(cCtx *cli.Context) error {

	fileName := cCtx.String("file")
	dnsName := strings.Split(fileName, ".")[0]
	domainName := cCtx.String("domain")

	log.Infof("Retrieving %s.", fileName)

	fileData, err := dnsToFile(dnsName, domainName)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("DNS-%s", fileName), fileData, os.ModeAppend)

	return err

}

func dnsToFile(filename string, domain string) ([]byte, error) {

	logger := log.StandardLogger().WithFields(log.Fields{
		"filename": filename,
		"domain":   domain,
	})

	sfile := ""
	// TODO: Figure out how to calculate this automatically
	for i := 0; i < 7; i++ {
		lookup := fmt.Sprintf("%s-%d.%s", filename, i, domain)
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
