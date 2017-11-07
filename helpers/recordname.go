package helpers

import (
	"strings"
)

func GenerateRecordName(hostname string, zone string, relative bool) string {
	recordName := strings.TrimRight(hostname, ".")
	if !relative {
		recordName = strings.TrimSuffix(hostname, zone)
		recordName = strings.TrimRight(recordName, ".")
	}

	if recordName == "" {
		return "@"
	}

	return recordName
}
