package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/elyscape/az-dns/helpers"
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
currently-supported record types are A, AAAA, and TXT. HOSTNAME may be a
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
        600 pointing at ::1`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		recordType := dns.RecordType(strings.ToUpper(args[0]))
		hostname := args[1]
		records := args[2:]

		client, err := helpers.NewRecordSetClient(dns.DefaultBaseURI)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		resourceGroup := viper.GetString("resource-group")
		if resourceGroup == "" {
			fmt.Println("A resource group name is required.")
			os.Exit(1)
		}

		zone := viper.GetString("zone")
		if zone == "" {
			fmt.Println("A DNS zone name is required.")
			os.Exit(1)
		}

		relative := viper.GetBool("relative")
		ttl := viper.GetInt64("ttl")
		recordName := helpers.GenerateRecordName(hostname, zone, relative)

		var rrparams *dns.RecordSet = nil
		switch recordType {
		case dns.A:
			rrparams, err = generateARecordParams(ttl, records)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case dns.AAAA:
			rrparams, err = generateAaaaRecordParams(ttl, records)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case dns.TXT:
			rrparams, err = generateTxtRecordParams(ttl, records)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unsupported record type %v.\n", recordType)
			os.Exit(1)
		}

		_, err = client.CreateOrUpdate(resourceGroup, zone, recordName, recordType, *rrparams, "", "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Success.")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.PersistentFlags().BoolP("relative", "r", false, "HOSTNAME is a zone-relative label")
	setCmd.PersistentFlags().Int64P("ttl", "t", 300, "Record set TTL")
	viper.BindPFlags(setCmd.PersistentFlags())
}

func generateARecordParams(ttl int64, values []string) (rrparams *dns.RecordSet, err error) {
	records := []dns.ARecord{}

	for _, addr := range values {
		if ip := net.ParseIP(addr); ip == nil || ip.To4() == nil {
			err = fmt.Errorf(`Invalid IP address "%v".`, addr)
			return
		}
		records = append(records, dns.ARecord{Ipv4Address: &addr})
	}

	rrparams = &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:      &ttl,
			ARecords: &records,
		},
	}

	return
}

func generateAaaaRecordParams(ttl int64, values []string) (rrparams *dns.RecordSet, err error) {
	records := []dns.AaaaRecord{}

	for _, addr := range values {
		if ip := net.ParseIP(addr); ip == nil || ip.To16() == nil {
			err = fmt.Errorf(`Invalid IP address "%v".`, addr)
			return
		}
		records = append(records, dns.AaaaRecord{Ipv6Address: &addr})
	}

	rrparams = &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:         &ttl,
			AaaaRecords: &records,
		},
	}

	return
}

func generateTxtRecordParams(ttl int64, values []string) (rrparams *dns.RecordSet, err error) {
	records := []dns.TxtRecord{}

	for _, value := range values {
		valueArr := []string{value}
		records = append(records, dns.TxtRecord{Value: &valueArr})
	}

	rrparams = &dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:        &ttl,
			TxtRecords: &records,
		},
	}

	return
}
