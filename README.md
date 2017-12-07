# az-dns
A simple command-line tool for manipulating [Azure DNS] record sets

This is a simple tool for managing Azure DNS resource record sets, written with
a primary focus on easy scriptability for use with things like Let's Encrypt DNS
challenge hooks. It authenticates to the service using an [Azure Active
Directory service principal][service principal], allowing you to limit which
records it can manage.

## Installation

Download the binary appropriate to your platform from [the GitHub releases
page][releases] and put it somewhere in your PATH. Alternatively, if you have Go
installed on your system, you can use that to install it:
```
go get -u github.com/elyscape/az-dns
```

## Configuration

This tool's various commands all need to know the Azure subscription ID, the
Azure resource group, and the name of the DNS zone on which they need to
operate. In addition to command-line flags, there are various ways that these
and other values can be provided.

### File

You can specify default values for flags by using a config file. On startup, the
tool will check your home directory for the existence of a file named
`.az-dns` with one of the following extensions and load the first one
that it finds:

- `.json`
- `.toml`
- `.yaml`
- `.yml`
- `.properties`
- `.props`
- `.prop`
- `.hcl`

The file must be a valid JSON, TOML, YAML, Java properties, or HCL document.
Valid keys are the long flag names (e.g. `resource-group`). Invalid fields are
ignored. You can specify a custom config file location by using the `config`
flag.

### Environment variables

With the exception of `config`, every flag can be specified through an
environment variable. The environment variable corresponding to a given flag is
`AZURE_`, followed by the name of the flag, in upper case, with all `-`
characters replaced with `_`s. For example, the environment variable
corresponding to `resource-group` is `AZURE_RESOURCE_GROUP`.

## Credentials

This tool needs the credentials for an Azure AD security principal in order to
authenticate to the Azure DNS servers.

### Generating credentials

To create a service principal and generate credentials, you will need to have
installed and configured the [Azure CLI]. Once it is installed and authenticated
to your account, you can create a service principal with permissions to modify
all DNS records on your account by running the following command:
```shellsession
$ az ad sp create-for-rbac --name AllDNS --role 'DNS Zone Contributor'
AppId                                 DisplayName    Name           Password                              Tenant
------------------------------------  -------------  -------------  ------------------------------------  ------------------------------------
12345678-90ab-cdef-1234-567890abcdef  AllDNS         http://AllDNS  fedcba09-8765-4321-fedc-ba0987654321  abcdef12-3456-7890-abcd-ef1234567890
```
The values for AppId, Password, and Tenant are the Azure client ID, client
secret, and tenant ID, respectively.

### Creating scoped credentials

To create a service principal that only has permissions to modify certain record
sets, you can specify a scope as well. Scopes look like this:
```
/subscriptions/<SUBSCRIPTION-UUID>/resourceGroups/dns/providers/Microsoft.Network/dnszones/<ZONE-NAME>/<RECORD-TYPE>/<LABEL>
```
For example, if you have a DNS zone for example.com and want to create a service
principal that has permissions to read and modify TXT records on
sub.example.com, the scope would be:
```
/subscriptions/<SUBSCRIPTION-UUID>/resourceGroups/dns/providers/Microsoft.Network/dnszones/example.com/TXT/sub
```
To allow access to the root of the domain (e.g. example.com), use `@` as the
label:
```
/subscriptions/<SUBSCRIPTION-UUID>/resourceGroups/dns/providers/Microsoft.Network/dnszones/example.com/A/@
```
One or more scopes can be provided to the Azure CLI using the `scopes` flag, like so:
```
az ad sp create-for-rbac --name SubdomainDNS --role 'DNS Zone Contributor' --scopes /subscriptions/<SUBSCRIPTION-UUID>/resourceGroups/dns/providers/Microsoft.Network/dnszones/example.com/TXT/sub /subscriptions/<SUBSCRIPTION-UUID>/resourceGroups/dns/providers/Microsoft.Network/dnszones/example.com/A/@
```

### Credential auth files

Instead of specifying the credentials through command line flags or in the
tool-specific config file, you can use an Azure SDK auth file. This file is
trivial for developers to support and will therefore likely be usable with other
tools without modification. To generate an Azure SDK auth file, pass the
`sdk-auth` flag to Azure CLI:
```shellsession
$ az ad sp create-for-rbac -n AuthFile --sdk-auth
{
  "clientId": "[REMOVED]",
  "clientSecret": "[REMOVED]",
  "subscriptionId": "[REMOVED]",
  "tenantId": "[REMOVED]",
  "activeDirectoryEndpointUrl": "https://login.microsoftonline.com",
  "resourceManagerEndpointUrl": "https://management.azure.com/",
  "activeDirectoryGraphResourceId": "https://graph.windows.net/",
  "sqlManagementEndpointUrl": "https://management.core.windows.net:8443/",
  "galleryEndpointUrl": "https://gallery.azure.com/",
  "managementEndpointUrl": "https://management.core.windows.net/"
}
```
Save this output into a file somewhere. To instruct the tool to use it, simply
provide the path to the file in the environment variable `AZURE_AUTH_LOCATION`.

[Azure DNS]: https://azure.microsoft.com/en-us/services/dns/
[service principal]: https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-application-objects
[releases]: https://github.com/elyscape/az-dns/releases/latest
[Azure CLI]: https://docs.microsoft.com/en-us/cli/azure/overview
