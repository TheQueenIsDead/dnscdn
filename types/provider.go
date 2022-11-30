package types

type DnsProviderType int

const (
	ProviderCloudflare DnsProviderType = iota
)

type DnsProvider interface {
	// Authenticate searches the environment for pre-defined strings in order to verify a
	Authenticate() error
	Blockify(filename string, domain string, blocks []string) error
	CreateRecord(fqdn string, data string) error
	ReadRecord() error
	UpdateRecord(fqdn string, data string) error
	DeleteRecord(fqdn string) error
}
