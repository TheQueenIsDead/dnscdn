package commands

import (
	"dnscdn/config"
	"dnscdn/provider"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func UploadCommand(cCtx *cli.Context) error {

	fileName := cCtx.String("file")

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
	err = dnsProvider.Blockify(fileName, "tqid.dev", blocks)
	if err != nil {
		log.WithError(err).Error("Error encountered calling Blockify()")
		return nil
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
