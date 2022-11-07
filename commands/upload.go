package commands

import (
	"github.com/urfave/cli/v2"
)

func UploadCommand(cCtx *cli.Context) error {
	panic("Not implemented")
	return nil
}

//
//func UploadCommand(cCtx *cli.Context) error {
//
//	// Retrieve filepath from args
//	filepath := os.Args[1]
//	log.Println(filepath)
//
//	// Load file
//	file, err := os.ReadFile(filepath)
//	if err != nil {
//		log.Fatalf("Could not open provided file (%s): %s", filepath, err)
//	}
//
//	// Encode
//	b64File := base64.StdEncoding.EncodeToString(file)
//
//	// Calculate blocks
//	blocks := stringToBlocks(b64File)
//
//	summary(blocks)
//
//	dnsProvider := provider.CloudflareDnsProvider{}
//	err = dnsProvider.Authenticate()
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	err = dnsProvider.Blockify(filepath, "tqid.dev", blocks)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//}
//
//func stringToBlocks(s string) []string {
//
//	// Calculate the amount of blocks required
//	nBlocks := len(s) / DnsCdnDataSize
//	if (len(s) % DnsCdnDataSize) != 0 {
//		nBlocks++
//	}
//
//	blocks := make([]string, nBlocks)
//
//	blockIndex := 0
//	blockCounter := 0
//	for _, char := range s {
//		blocks[blockIndex] = blocks[blockIndex] + string(char)
//
//		blockCounter++
//
//		// If we have exceeded the max allowed block size, move to the next block string by
//		// incrementing the index
//		if blockCounter >= dnsCdnDataSize {
//			blockCounter = 0
//			blockIndex++
//		}
//		//log.Debugf("Rune %v is '%c' (%d)\n", blockIndex, char, blockCounter)
//	}
//
//	return blocks
//}
//
//func summary(blocks []string) {
//	for i, b := range blocks {
//		log.Debugf("Block %v, len: %v", i, len(b))
//	}
//}
