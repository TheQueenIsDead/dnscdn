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

	// Create index record
	idxName := fmt.Sprintf("%s.media.%s", filename, domain)
	idxData := len(blocks)
	err = c.createRecord(idxName, strconv.Itoa(idxData))
	if err != nil {
		log.WithError(err).Error("Could not create index record.")
		return err
	}

	// Create data records
	for i, b := range blocks {
		fqdn := fmt.Sprintf("%s.%v.media.%s", filename, i, domain)
		err := c.createRecord(fqdn, b)
		if err != nil {
			log.Fatalf("Failed to create block %d of %s: %v", i, filename, err)
			return err
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
			log.WithField("fqdn", fqdn).Warn("Duplicate record present removing...")
			c.DeleteRecord(fqdn)
			c.createRecord(fqdn, data)
		}
		log.Error(err.Error())

	} else {
		log.Debugf("Successfully created record: %s", fqdn)
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

func (c *CloudflareDnsProvider) DeleteRecord(fqdn string) error {

	log.Warnf("Deleting Cloudflare DNS records for %s", fqdn)

	// Fetch records of any type with name "foo.example.com"
	// The name must be fully-qualified
	foo := cloudflare.DNSRecord{Name: fqdn}
	recs, err := c.api.DNSRecords(context.Background(), c.apiZone, foo)
	if err != nil {
		log.Fatal(err)
	}

	// Enumerate and delete
	for _, rec := range recs {
		err := c.api.DeleteDNSRecord(context.Background(), c.apiZone, rec.ID)
		if err != nil {
			return err
		}
	}

	return err
}
