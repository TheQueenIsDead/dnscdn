package provider

import (
	"context"
	"dnscdn/lib"
	"errors"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"os"
)

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

	return err
}

func (c *CloudflareDnsProvider) Blockify(fileName string, domainName string, blocks []string) error {

	var err error

	// Create data records
	for i, b := range blocks {
		fqdn := lib.DataFqdn(fileName, domainName, i)
		err := c.CreateRecord(fqdn, b)
		if err != nil {
			log.WithError(err).Error("Failed to create block %d of %s: %v", i, fileName, err)
			return err
		}
	}

	return err
}

func (c *CloudflareDnsProvider) CreateRecord(fqdn string, data string) error {

	_, err := c.api.CreateDNSRecord(context.Background(), c.apiZone, cloudflare.DNSRecord{
		Type:    "TXT",
		Name:    fqdn,
		Content: data,
		ZoneID:  c.apiZone,
	})

	if err != nil {
		log.WithError(err).Error("Could not create record")
	} else {
		log.Debugf("Successfully created record: %s", fqdn)
	}

	return err
}

func (c *CloudflareDnsProvider) ReadRecord() error {

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

func (c *CloudflareDnsProvider) UpdateRecord(fqdn string, data string) error {

	logger := log.StandardLogger().WithFields(log.Fields{
		"fqdn": fqdn,
	})

	// Fetch records of any type with name "foo.example.com"
	// The name must be fully-qualified
	foo := cloudflare.DNSRecord{Name: fqdn}
	recs, err := c.api.DNSRecords(context.Background(), c.apiZone, foo)
	if err != nil {
		logger.Fatal(err)
	}

	if len(recs) != 1 {
		return errors.New("more than one record returned while trying to update")
	}

	err = c.api.UpdateDNSRecord(context.Background(), c.apiZone, recs[0].ID, cloudflare.DNSRecord{
		Type:    "TXT",
		Name:    fqdn,
		Content: data,
		ZoneID:  c.apiZone,
	})

	if err != nil {
		logger.WithError(err).Error("Could not update record.")
	} else {
		logger.Info("Record updated.")
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
