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

	filename := cli.
		log.Info("Retrieving file!")
	dnsFile := dnsToFile(filepath)
	err := os.WriteFile(fmt.Sprintf("reconstructed-%s", filepath), dnsFile, os.ModeAppend)
	if err != nil {
		return
	}
}

func stringToBlocks(s string) []string {

	// Calculate the amount of blocks required
	nBlocks := len(s) / DnsCdnDataSize
	if (len(s) % DnsCdnDataSize) != 0 {
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
		if blockCounter >= dnsCdnDataSize {
			blockCounter = 0
			blockIndex++
		}
		//log.Debugf("Rune %v is '%c' (%d)\n", blockIndex, char, blockCounter)
	}

	return blocks
}

func dnsToFile(filename string) []byte {

	// TODO: Pass in the domain
	domain := "tqid.dev"
	filename = strings.Split(filename, ".")[0]

	sfile := ""
	for i := 0; i < 7; i++ {
		lookup := fmt.Sprintf("%s-%d.%s", filename, i, domain)
		log.Debugf("Retreiving TXT record for %s", lookup)
		txt, _ := net.LookupTXT(lookup)
		sfile = sfile + txt[0]
	}

	file, _ := base64.StdEncoding.DecodeString(sfile)

	return file

}

func summary(blocks []string) {
	for i, b := range blocks {
		log.Debugf("Block %v, len: %v", i, len(b))
	}
}
