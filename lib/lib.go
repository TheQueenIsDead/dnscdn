package lib

import (
	"dnscdn/config"
	types "dnscdn/types"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// IndexFqdn returns a standard index FQDN given a domain.
// Returns a string in the format dnscdn.<domain>
func IndexFqdn(domainName string) string {
	return fmt.Sprintf("%s.%s", config.DnsCdnIndexDomain, domainName)
}

// DataFqdn returns a standard FQDN given a filename, block number and domain.
// Returns a string in the format <file.name>.<index>.dnscdn.<domain>
// Eg, "melvin.png.3.dnscdn.tqid.dev"
func DataFqdn(fileName string, domainName string, index int) string {
	return fmt.Sprintf("%s.%v.%s.%s", fileName, index, config.DnsCdnIndexDomain, domainName)
}

// ParseIndex takes a string from an index TXT record and parses it into a map of
// filenames to their corresponding record count.
func ParseIndex(data string) map[string]int {
	records := strings.Split(data, config.DnsCdnIndexSeparator)

	results := map[string]int{}

	for _, r := range records {
		components := strings.Split(r, ",") // filename,length

		filename := components[0]
		length, err := strconv.Atoi(components[1])
		if err != nil {
			log.WithError(err).Error("Could not parse length of record in index.")
		}

		results[filename] = length
	}

	return results
}

// UpdateIndex takes a map of filenames to their corresponding record count and forms a string which
// then updates an existing index TXT record, or creates on if the update fails.
// If the string is empty, the record is deleted.
// The string is a concatenation of <filename>,<length> records separated by semicolons
// Ex, melvin.png,7;nivlem.gnp,1
func UpdateIndex(domainName string, records map[string]int, provider types.DnsProvider) error {

	// Parse index string
	idxFqdn := IndexFqdn(domainName)

	// Build data string based on records
	data := ""
	for filename, length := range records {
		data += strings.Join([]string{filename, strconv.Itoa(length)}, ",") + ";"
	}
	data = strings.TrimSuffix(data, ";")

	// If the string is empty attempt to delete the index record
	if data == "" {
		return provider.DeleteRecord(idxFqdn)
	}

	// Otherwise, try to update the index
	err := provider.UpdateRecord(idxFqdn, data)
	if err != nil {
		// Hail mary create, assuming the error is that it does not exist
		err := provider.CreateRecord(idxFqdn, data)
		if err != nil {
			return err
		}
	}

	return err
}
