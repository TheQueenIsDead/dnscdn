package lib

import (
	"dnscdn/config"
	"fmt"
)

// IndexFqdn returns a standard index FQDN given a domain.
// Returns a string in the format dnscdn.<domain>
func IndexFqdn(domainName string) string {
	return fmt.Sprintf("%s.%s", config.DnsCdnIndexDomain, domainName)
}

// DataFqdn returns a standard FQDN given a filename, block number and domain.
// Returns a string in the format <file.name>.<index>.media.<domain>
func DataFqdn(fileName string, index int, domainName string) string {
	return fmt.Sprintf("%s.%v.media.%s", fileName, index, domainName)
}
