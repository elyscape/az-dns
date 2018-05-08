package helpers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	authfile "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/viper"
)

// NewRecordSetClient creates a new RecordSetsClient using the specified
// baseURI and attaches a BearerAuthorizer based on credentials provided either
// via an Azure SDK auth file, if present, or through any mechanism supported
// by Viper. If credentials have not been provided, an error will be returned.
func NewRecordSetClient(baseURI string) (*dns.RecordSetsClient, error) {
	var authorizer *autorest.BearerAuthorizer
	subscriptionID := viper.GetString("subscription-id")

	if clientSetup, err := authfile.GetClientSetup(baseURI); err == nil {
		authorizer = clientSetup.BearerAuthorizer
		subscriptionID = clientSetup.SubscriptionID
	} else {
		authorizer, err = GetAuthorizer(baseURI)
		if err != nil {
			return nil, err
		}
	}

	client := dns.NewRecordSetsClientWithBaseURI(baseURI, subscriptionID)
	client.Authorizer = authorizer

	return &client, nil
}

// GetAuthorizer creates a BearerAuthorizer based on credentials retrieved from
// Viper. If credentials have not been provided, an error will be returned.
func GetAuthorizer(baseURI string) (*autorest.BearerAuthorizer, error) {
	credFields := []string{"client-id", "client-secret", "subscription-id", "tenant-id"}

	for _, field := range credFields {
		if !viper.IsSet(field) || viper.GetString(field) == "" {
			return nil, fmt.Errorf("required credential option %v not provided", field)
		}
	}

	config, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, viper.GetString("tenant-id"))
	if err != nil {
		return nil, err
	}

	token, err := adal.NewServicePrincipalToken(*config, viper.GetString("client-id"), viper.GetString("client-secret"), baseURI)
	if err != nil {
		return nil, err
	}

	return autorest.NewBearerAuthorizer(token), nil
}
