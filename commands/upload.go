package commands

import (
	"dnscdn/provider"
	"encoding/base64"
	"github.com/urfave/cli/v2"
	"os"
)

func UploadCommand(cCtx *cli.Context) error {

	// Retrieve filepath from args
	filepath := os.Args[1]
	log.Println(filepath)

	// Load file
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Could not open provided file (%s): %s", filepath, err)
	}

	// Encode
	b64File := base64.StdEncoding.EncodeToString(file)

	// Calculate blocks
	blocks := stringToBlocks(b64File)

	summary(blocks)

	dnsProvider := provider.CloudflareDnsProvider{}
	err = dnsProvider.Authenticate()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = dnsProvider.Blockify(filepath, "tqid.dev", blocks)
	if err != nil {
		log.Fatal(err)
		return
	}

}
