package config

import log "github.com/sirupsen/logrus"

const (
	// DnsCdnDataSize TODO: If this shrinks then it is possible that new TXT records will be created as the data is different.
	DnsCdnDataSize        = 2048
	DnsCdnIndexSeparator  = ";"
	DsnCdnDefaultLogLevel = log.DebugLevel
	DnsCdnIndexDomain     = "dnscdn"
)
