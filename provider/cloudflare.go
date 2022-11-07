package provider

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

// TODO: If this shrinks then it is possible that new TXT records will be created as the data is different.
const ErrRecordAlreadyExistsCode = 81057

type CloudflareDnsProvider struct {
	api      cloudflare.API
	apiKey   string
	apiEmail string
	apiZone  string
}

func (c *CloudflareDnsProvider) Authenticate() error {

	c.apiKey = os.Getenv("CLOUDFLARE_API_KEY")
	c.apiEmail = os.Getenv("CLOUDFLARE_API_EMAIL")
	c.apiZone = os.Getenv("CLOUDFLARE_API_ZONE")

	if c.apiKey == "" || c.apiEmail == "" || c.apiZone == "" {
		log.Fatal("Could not retrieve Cloudflare configuration options from the environment.")
	}

	api, err := cloudflare.New(c.apiKey, c.apiEmail)
	if err != nil {
		log.Fatalf("Error authenticating with Cloudflare: %v", err)
		return err
	}
	log.Infof("Successfully authenticated with Cloudflare.")
	c.api = *api

	return nil
}

func (c *CloudflareDnsProvider) Blockify(filename string, domain string, blocks []string) error {

	var err error

	filename = strings.Split(filename, ".")[0]

	for i, b := range blocks {
		fqdn := fmt.Sprintf("%s-%v", filename, i)
		//fqdn := fmt.Sprintf("%s-%v.%s", filename, i, domain)
		err := c.createRecord(fqdn, b)
		if err != nil {
			log.Fatalf("Failed to create block %d of %s: %v", i, filename, err)
			break
		} else {
			log.Debugf("Successfully created record: %s", fqdn)
		}
	}

	return err
}

func (c *CloudflareDnsProvider) createRecord(fqdn string, data string) error {

	_, err := c.api.CreateDNSRecord(context.Background(), c.apiZone, cloudflare.DNSRecord{
		Type:    "TXT",
		Name:    fqdn,
		Content: data,
		ZoneID:  c.apiZone,
	})

	if err != nil {

		serr := err.Error()
		if strings.Contains(serr, strconv.Itoa(ErrRecordAlreadyExistsCode)) {
			log.Warn("Duplicate record present removing...")
			c.deleteRecord(fqdn)
			c.createRecord(fqdn, data)
		}
		log.Error(err.Error())

		//switch errc {
		//default:
		//	log.Errorf("Unknown error: %v", errc)
		//case ErrRecordAlreadyExistsCode:
		//	log.Error("Record already exists tho")
		//}

	}

	return err
}

func (c *CloudflareDnsProvider) readRecord() error {

	// Fetch all records for a zone
	recs, err := c.api.DNSRecords(context.Background(), c.apiZone, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range recs {
		fmt.Printf("%s: %s\n", r.Name, r.Content)
	}

	return err
}

func (c *CloudflareDnsProvider) deleteRecord(fqdn string) error {

	// TODO: Pass in the full domain to provider
	fqdn = fqdn + ".tqid.dev"

	log.Warnf("Deleting Cloudflare DNS records for %s", fqdn)

	// Fetch records of any type with name "foo.example.com"
	// The name must be fully-qualified
	foo := cloudflare.DNSRecord{Name: fqdn}
	recs, err := c.api.DNSRecords(context.Background(), c.apiZone, foo)
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range recs {
		log.Infof("Deleting Cloudflare DNS record %d", rec.ID)
		err := c.api.DeleteDNSRecord(context.Background(), c.apiZone, rec.ID)
		if err != nil {
			return err
		}
	}

	return err
}
