package provider

type DnsProviderType int

const (
	ProviderCloudflare DnsProviderType = iota
)

type DnsProvider interface {
	Authenticate() error
	Blockify(filename string, domain string, blocks []string) error
	CreateRecord(fqdn string, data string) error
	ReadRecord() error
	UpdateRecord() error
	DeleteRecord(fqdn string) error
}
