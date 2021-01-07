# GCloud Function POC

This is a proof of concept to demonstrate how serverless computing can be supported in SPIRE through the introduction of an `SVIDStore` agent plugin.

The model leverages the use of secret management services offered by cloud providers to store and retrieve the SVIDs and keys in a secure way, inside the cloud infrastructure.

The serverless functions are registered in SPIRE in the same way that regular workloads are registered through registration entries. The `svidstore` key is used to distinguish the "storable" entries, and `SVIDStore` plugins receive updates of those entries only, which indicates that the issued SVID and key must be securely stored in a location accessible by the serverless function. This way, selectors provide a flexible way to describe the attributes needed to store the corresponding issued SVID and key, like the type of store, name to provide to the secret, and any specific attribute needed by the specific service used.

## Components

### AWS Lambda Extension

Simple extension that reads a secret from AWS Secrets Manager.

*NOTE: it is expected that the secret is a binary `workload.X509SVIDResponse` message.*

The secret name or ARN must be provided using the environment variable "SECRET_NAME" in the function. 
The secret is parsed and the X509-SVID, bundle and key are persisted in the `/tmp` folder.

### Function

Simple function that acces to SecretManger API, and returns a X509-SVID stored in a Secret.

## Scripts

* [00-var.sh](./00-vars.sh): Contains all the variables used to run this POC. 'PROJECT_ID' and 'SERVICE_ACCOUNT' must be updated.
* [01-deploy-function.sh](./01-deploy-functions.sh): Deploy function.
* [02-get.sh](./02-get.sh): Test function running a GET

## SPIRE changes

The SPIRE Agent cache manager was updated to be able to identify "storable" entries and notify the corresponding plugin when the entries are updated. A new `SVIDStore` agent plugin is introduced for this.

### Entry example

```
Entry ID         : 15a09d11-85b7-4ecb-8816-ac1b6f293d03
SPIFFE ID        : spiffe://example.org/gcloud/poc-function
Parent ID        : spiffe://example.org/agent
TTL              : default 
Selector         : svidstore:secretname:svid-dbuser
Selector         : svidstore:secretproject:sample-project
Selector         : svidstore:type:gcloud_secretsmanager
```

* `svidstore` is the key used to indicate that the issued SVID and key must be stored in a secure store.
* `svidstore:type:gcloud_secretsmanager` indicates that the entry must be stored in a store of type `gcloud_secretsmanager`.
* `svidstore:secretname:svid-dbuser` is an example of how a platform-specific attribute can be specified. in this case it indicates Secret name.
* `svidstore:secretproject:sample-project` project name where secret is created.

## Configuration sample:

```
SVIDStore "gcloud_secretsmanager" {
    plugin_data {
    	service_account_file = "/Users/marcosyacob/Downloads/dev-myacob-c24253d99f22.json"
    }
}
```

* `service_account_file` contains credentials used when connecting to Secrets Manager API, this SA must has permissions to access "secretmanager.secrets.get", "secretmanager.versions.add"

