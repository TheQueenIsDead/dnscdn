package provider

type DnsProviderType int

const (
	ProviderCloudflare DnsProviderType = iota
)

type DnsProvider interface {
	Authenticate() error
	Blockify(filename string, blocks []string) error
	createRecord() error
	readRecord() error
	deleteRecord() error
}
