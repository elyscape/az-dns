package helpers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
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
func NewRecordSetClient(baseURI string) (client dns.RecordSetsClient, err error) {
	var authorizer *autorest.BearerAuthorizer
	subscriptionID := viper.GetString("subscription-id")

	clientSetup, err := authfile.GetClientSetup(baseURI)
	if err == nil {
		authorizer = clientSetup.BearerAuthorizer
		subscriptionID = clientSetup.SubscriptionID
	} else {
		authorizer, err = GetAuthorizer(baseURI)
		if err != nil {
			return
		}
	}

	client = dns.NewRecordSetsClientWithBaseURI(baseURI, subscriptionID)
	client.Authorizer = authorizer

	return
}

// GetAuthorizer creates a BearerAuthorizer based on credentials retrieved from
// Viper. If credentials have not been provided, an error will be returned.
func GetAuthorizer(baseURI string) (authorizer *autorest.BearerAuthorizer, err error) {
	credFields := []string{"client-id", "client-secret", "subscription-id", "tenant-id"}

	for _, field := range credFields {
		if !viper.IsSet(field) || viper.GetString(field) == "" {
			err = fmt.Errorf("required credential option %v not provided", field)
			return
		}
	}

	config, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, viper.GetString("tenant-id"))
	if err != nil {
		return
	}

	token, err := adal.NewServicePrincipalToken(*config, viper.GetString("client-id"), viper.GetString("client-secret"), baseURI)
	if err != nil {
		return
	}

	authorizer = autorest.NewBearerAuthorizer(token)
	return
}
