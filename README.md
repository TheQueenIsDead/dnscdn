# DNSCDN

A Golang CLI for storing files on the great Domain Name System. Who doesn't need one more reason to [debug DNS](https://isitdns.com/)?

## Overview

I started wondering what limits are on DNS records, in terms of content size, which led me down the road of wondering if 
it would be terribly difficult to store data in chunks across a variety of records.

As it happens it's rather trivial. More about this project can be read on my [blog](https://tqid.dev/posts/dnscdn/)

## Usage

Download
```console
dnscdn download -f melvin.png -d tqid.dev 
```

Upload
```console
dnscdn upload -f melvin.png -d tqid.dev 
```

List
```console
dnscdn list -f tqid.dev

DNSCDN files on tqid.dev:
Filename             Length     Size (kB) 
melvin.png           7          14       
```

Delete
```console
dnscdn delete -f melvin.png -d tqid.dev
```

Help
```console
NAME:
   DNSCDN - Store and retrieve media by use of 'free' DNS storage.

USAGE:
   DNSCDN [global options] command [command options] [arguments...]

COMMANDS:
   upload    Upload a file to a given DNS provider by means of TXT record.
   download  Retrieve file data from DNS and save it locally.
   delete    Remove all file specific records for a domain.
   list      Given a domain, enumerate for previously saved DNSCDN media.
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```
