# DNSCDN

A Golang CLI for storing files on the great Domain Name System. Who doesn't need one more reason to [debug DNS](https://isitdns.com/)?

## Overview

I started wondering what limits are on DNS records, in terms of content size, which led me down the road of wondering if 
it would be terribly difficult to store data in chunks across a variety of records.

As it happens it's rather trivial. More about this project can be read on my [blog](https://tqid.dev/posts/dnscdn/)

## Usage

Download
```console
dnscdn -f melvin.png -d tqid.dev download
```

Upload
```console
dnscdn -f melvin.png -d tqid.dev upload
```

List
```console
dnscdn -f tqid.dev list

DNSCDN files on tqid.dev:
Filename             Length     Size (kB) 
melvin.png           7          14       
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
