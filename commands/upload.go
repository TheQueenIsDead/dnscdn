package commands

import (
	"dnscdn/config"
	"dnscdn/lib"
	"dnscdn/provider"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"os"
)

func UploadCommand(cCtx *cli.Context) error {

	fileName := cCtx.String("file")
	domainName := cCtx.String("domain")
	idxFqdn := lib.IndexFqdn(domainName)

	// Lookup index record in order to see if file already exists
	records := map[string]int{}
	idx, err := net.LookupTXT(idxFqdn)
	if err == nil {
		// Exit if the index already contains the file we're trying to retrieve
		if _, ok := records[fileName]; ok {
			return errors.New(fmt.Sprintf("File '%s' already exists in domain index.", fileName))
		}
		records = lib.ParseIndex(idx[0])
	}

	// Load file
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Errorf("Could not open provided file (%s): %s", fileName, err)
		return err
	}

	// Encode
	b64File := base64.StdEncoding.EncodeToString(file)

	// Calculate blocks
	blocks := stringToBlocks(b64File, config.DnsCdnDataSize)

	summary(blocks)

	dnsProvider := provider.CloudflareDnsProvider{}
	err = dnsProvider.Authenticate()
	if err != nil {
		log.WithError(err).Error("Error encountered calling Authenticate()")
		return nil
	}
	err = dnsProvider.Blockify(fileName, domainName, blocks)
	if err != nil {
		log.WithError(err).Error("Error encountered calling Blockify()")
		return nil
	}

	// Update index
	records[fileName] = len(blocks)
	err = lib.UpdateIndex(domainName, records, &dnsProvider)
	if err != nil {
		return err
	}

	return nil
}

func stringToBlocks(s string, blockSize int) []string {

	// Calculate the amount of blocks required
	nBlocks := len(s) / blockSize
	if (len(s) % blockSize) != 0 {
		nBlocks++
	}

	blocks := make([]string, nBlocks)

	blockIndex := 0
	blockCounter := 0
	for _, char := range s {
		blocks[blockIndex] = blocks[blockIndex] + string(char)

		blockCounter++

		// If we have exceeded the max allowed block size, move to the next block string by
		// incrementing the index
		if blockCounter >= blockSize {
			blockCounter = 0
			blockIndex++
		}
		//log.Debugf("Rune %v is '%c' (%d)\n", blockIndex, char, blockCounter)
	}

	return blocks
}

func summary(blocks []string) {
	for i, b := range blocks {
		log.Debugf("Block %v, len: %v", i, len(b))
	}
}
