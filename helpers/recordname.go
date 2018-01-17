package helpers

import (
	"strings"
)

// GenerateRecordName determines a DNS resource record name from a hostname and
// the zone name. It strips any trailing dots and optionally the zone name from
// the provided hostname. If the result would be the empty string, it instead
// returns "@", which corresponds to the apex of DNS zone.
func GenerateRecordName(hostname string, zone string, relative bool) string {
	recordName := strings.TrimRight(hostname, ".")
	if !relative {
		recordName = strings.TrimSuffix(recordName, zone)
		recordName = strings.TrimRight(recordName, ".")
	}

	if recordName == "" {
		return "@"
	}

	return recordName
}
