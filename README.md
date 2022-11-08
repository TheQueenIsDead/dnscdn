# DNSCDN

A Golang CLI for storing files on the great Domain Name System. Who doesn't need one more reason to [debug DNS](https://isitdns.com/)?

## Overview

I started wondering what limits are on DNS records, in terms of content size, which led me down the road of wondering if 
it would be terribly difficult to store data in chunks across a variety of records.

As it happens it's rather trivial.

## Usage

Download
```console
dnscdn -f melvin.png -d tqid.dev download
```

Upload
```console
dnscdn -f melvin.png -d tqid.dev upload
```

Delete
```console
dnscdn -f melvin.png -d tqid.dev delete
```

Help
```console
dnscdn --help
NAME:
   DNSCDN - Store and retrieve media by use of 'free' DNS storage.

USAGE:
   DNSCDN [global options] command [command options] [arguments...]

COMMANDS:
   upload    Upload a file to a given DNS provider by means of TXT record.
   download  Retrieve file data from DNS and save it locally.
   delete    Remove all file specific records for a domain.
   list, d   Given a domain, enumerate for previously saved DNSCDN media.
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value    File to retrieve or upload, including extension.
   --domain value, -d value  Domain to retrieve from or upload to.
   --help, -h                show help (default: false)
```

## DNS

There's a few [DNS record types](https://en.wikipedia.org/wiki/List_of_DNS_record_types#Resource_records) out there, 
with some decent resources that shed light on the make-up of packets. 

 - https://aws.amazon.com/premiumsupport/knowledge-center/route-53-configure-long-spf-txt-records/
 - https://stackoverflow.com/questions/6794926/how-many-a-records-can-fit-in-a-single-dns-response
 - https://datatracker.ietf.org/doc/html/draft-ietf-dnsop-respsize-12
 - https://www.rfc-editor.org/rfc/rfc1035
 - https://blog.cloudflare.com/a-deep-dive-into-dns-packet-sizes-why-smaller-packet-sizes-keep-the-internet-safe/

As to the max? It quite depends on the provider.

 - AWS Route 53: [4,000 characters](The maximum length of a value in a TXT record is 4,000 characters.)
 - Cloudflare: [2048 characters](https://developers.cloudflare.com/dns/manage-dns-records/reference/dns-record-types/#txt)
 - Azure DNS: [1024 characters](https://learn.microsoft.com/en-us/azure/dns/dns-zones-records#txt-records)

These will each have different limits on the amount of records you can employ per domain, and how many of each type 
you may provision.

[//]: # (TODO: Create a table showing account limits for theoretical file sizes.)

I started with Cloudflare because they [do not charge additional provider specific fees](https://www.cloudflare.com/products/registrar/)
for management, and have a decent API to interface with.

## TXT

Excerpt from [AWS](https://aws.amazon.com/premiumsupport/knowledge-center/route-53-configure-long-spf-txt-records/#:~:text=The%20maximum%20length%20of%20a,TXT%20record%20is%204%2C000%20characters)

Key points to remember:

 - A TXT record contains one or more strings that are enclosed in double quotation marks (").
 - [AWS] You can enter a value of up to 255 characters in one string in a TXT record.
 - You can add multiple strings of 255 characters in a single TXT record.
 - The maximum length of a value in a TXT record is 4,000 characters.
 - TXT record values are case-sensitive.

For values that exceed 255 characters, break the value into strings of 255 characters or less. 
Enclose each string in double quotation marks (") using the following syntax: 
Domain name TXT "String 1" "String 2" "String 3"….."String N".

N.B: Cloudflare will do this for you
