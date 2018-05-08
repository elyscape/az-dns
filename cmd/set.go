package cmd

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/elyscape/az-dns/helpers"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set TYPE HOSTNAME VALUES",
	Short: "Create or update a DNS record set",
	Long: `Create or update a record set in Azure DNS

This will create or update a record set on Azure DNS, depending on whether a
record of the same type already exists for the provided value of HOSTNAME. The
currently-supported record types are A, AAAA, CAA, and TXT. HOSTNAME may be a
fully-qualified domain name contained within the zone, a record name relative
to the zone, or either the empty string or @ for the apex. If a record name
contains the zone name (e.g. example.com.example.com), you should either
provide the FQDN or use the --relative flag.

Examples:
    az-dns set A example.com 1.1.1.1 -z example.com
        Creates an A record at the apex of example.com pointing to 1.1.1.1
    az-dns set A sub 1.1.1.1 2.2.2.2 -z example.com
        Creates an A record for sub.example.com pointing to 1.1.1.1 and 2.2.2.2
    az-dns set AAAA local.example.com ::1 -t 600 -r -z example.com
        Creates an AAAA record for local.example.com.example.com with TTL of
        600 pointing at ::1
    az-dns set CAA example.com 0 issue letsencrypt.org -z example.com
        Creates a CAA record at the apex of example.com with value:
            0 issue "letsencrypt.org"
    az-dns set CAA @ 0 issue letsencrypt.org 0 issuewild ';' -z example.com
        Creates CAA records at the apex of example.com with values:
            0 issue "letsencrypt.org"
            0 issuewild ";"`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		recordType := dns.RecordType(strings.ToUpper(args[0]))
		hostname := args[1]
		records := args[2:]

		client, err := helpers.NewRecordSetClient(dns.DefaultBaseURI)
		if err != nil {
			return err
		}

		resourceGroup := viper.GetString("resource-group")
		if resourceGroup == "" {
			return fmt.Errorf("a resource group name is required")
		}

		zone := viper.GetString("zone")
		if zone == "" {
			return fmt.Errorf("a DNS zone name is required")
		}

		relative := viper.GetBool("relative")
		ttl := viper.GetInt64("ttl")
		recordName := helpers.GenerateRecordName(hostname, zone, relative)

		var rrparams *dns.RecordSet = nil
		switch recordType {
		case dns.A:
			rrparams, err = generateARecordParams(ttl, records)
			if err != nil {
				return err
			}
		case dns.AAAA:
			rrparams, err = generateAaaaRecordParams(ttl, records)
			if err != nil {
				return err
			}
		case dns.CAA:
			rrparams, err = generateCaaRecordParams(ttl, records)
			if err != nil {
				return err
			}
		case dns.TXT:
			rrparams, err = generateTxtRecordParams(ttl, records)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported record type %v", recordType)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err = client.CreateOrUpdate(ctx, resourceGroup, zone, recordName, recordType, *rrparams, "", "")
		if err != nil {
			return err
		}

		fmt.Println("success")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.PersistentFlags().BoolP("relative", "r", false, "HOSTNAME is a zone-relative label")
	setCmd.PersistentFlags().Int64P("ttl", "t", 300, "Record set TTL")
	if err := viper.BindPFlags(setCmd.PersistentFlags()); err != nil {
		// This shouldn't happen
		panic(err)
	}
}

func generateARecordParams(ttl int64, values []string) (*dns.RecordSet, error) {
	records := []dns.ARecord{}

	for _, addr := range values {
		if ip := net.ParseIP(addr); ip == nil || ip.To4() == nil {
			return nil, fmt.Errorf(`invalid IP address "%v"`, addr)
		}
		records = append(records, dns.ARecord{Ipv4Address: &addr})
	}

	rrparams := &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:      &ttl,
			ARecords: &records,
		},
	}

	return rrparams, nil
}

func generateAaaaRecordParams(ttl int64, values []string) (*dns.RecordSet, error) {
	records := []dns.AaaaRecord{}

	for _, addr := range values {
		if ip := net.ParseIP(addr); ip == nil || ip.To16() == nil {
			return nil, fmt.Errorf(`invalid IP address "%v"`, addr)
		}
		records = append(records, dns.AaaaRecord{Ipv6Address: &addr})
	}

	rrparams := &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:         &ttl,
			AaaaRecords: &records,
		},
	}

	return rrparams, nil
}

func generateCaaRecordParams(ttl int64, values []string) (*dns.RecordSet, error) {
	records := []dns.CaaRecord{}

	const recordSize = 3

	for min := 0; min < len(values); min += recordSize {
		max := min + recordSize
		if max > len(values) {
			return nil, fmt.Errorf(`incomplete CAA record %v`, values[min:])
		}

		fields := values[min:max]

		var flags int32
		flags, err := cast.ToInt32E(fields[0])
		if err != nil || flags > 255 || flags < 0 {
			return nil, fmt.Errorf(`invalid CAA flags "%v" must be an integer between 0 and 255`, fields[0])
		}

		tag := fields[1]
		value := fields[2]

		records = append(records, dns.CaaRecord{
			Flags: &flags,
			Tag:   &tag,
			Value: &value,
		})
	}

	rrparams := &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:        &ttl,
			CaaRecords: &records,
		},
	}

	return rrparams, nil
}

func generateTxtRecordParams(ttl int64, values []string) (*dns.RecordSet, error) {
	records := []dns.TxtRecord{}

	for _, value := range values {
		valueArr := []string{value}
		records = append(records, dns.TxtRecord{Value: &valueArr})
	}

	rrparams := &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:        &ttl,
			TxtRecords: &records,
		},
	}

	return rrparams, nil
}
