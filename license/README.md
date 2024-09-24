# UniPDF License Loading.

The examples here illustrate how to work with UniPDF license codes and keys.
There are two types of licenses.

## Offline License
Offline licenses are cryptography based and contain full signed information that is verified based on signatures without making any outbound connections,
hence the name "offline". This kind of license is suitable for users deploying OEM products to their customers or where there are strict restrictions
on outbound connections due to firewalls and/or compliance requirements.

## Metered License (API keys)
The metered license is the most convenient way to get started with UniDoc products and the Free tier enables a powerful way to get started for free.
Anyone can get a free metered API key by signing up on http://cloud.unidoc.io/

> Metered License (API keys) requires read-write permission to $HOME directory for storing API Keys usage,
> however you can set `HOME` environment to another directory if you wish.

## Metered License API Key Usage Logs
By setting the `SetMeteredKeyUsageLogVerboseMode` to true using `license.SetMeteredKeyUsageLogVerboseMode(true)` you can see full information on the credit usage status of each logs of each document as follows.
![alt text](data/image.png)   

## Examples

- [unipdf_license_loading_metered.go](unipdf_license_loading_metered.go) Demonstrates how to load the Metered API license key and how to print out relevant information.
- [unipdf_offline_license_info.go](unipdf_offline_license_info.go) Demonstrates how to print out information about the license after loading an offline license key.
- [unipdf_license_usage_log.go](unipdf_license_usage_log.go) Demonstrates how to enable the license key verbose mode logging.