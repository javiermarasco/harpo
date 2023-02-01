# Harpo (short for harpocrates)

The objective of this cli is to be able to manage secrets in Azure, AWS and GCP by using a path to 
specify where a secret is stored in a "folder structure" approach while providing an easy way to list 
secrets in this folder structure and also an easy way to retrieve the values of those secrets.


# How to build it

To build this project you will need Go 1.17 and run the following commands
```
go build
```

That will create a file called `harpo.exe` in Windows and `harpo` in Linux
# How to use

The syntax of the commands is the same for all cloud providers, just keep in mind the order of the parameters is important (currently is not possible to exchange the order of the parameters.)
## Environment variables setup

You need to define some environment variables in order to use `harpo`, those depend on the cloud provider and you can have more than one cloud provider setup at the same time (specially if you want to copy or migrate from one cloud to another this is mandatory). Once your have this variables defined you can start using the CLI.

### Azure
- "AZ_CLIENTID" (Contains the client id of the service princial/app registration used to access your keyvault)
- "AZ_CLIENTSECRET" (Contains the secret of the service principal/app registration)
- "AZ_TENANTID" (Contains the Tenant ID where your keyvault is deployed)
- "AZ_KVNAME" (Contains the name of the keyvault to use)

### AWS
- "AWS_ACCESS_KEY_ID" (Contains the key ID of the user that will be used to access Secrets Manager)
- "AWS_SECRET_ACCESS_KEY" (The access key of the user that will be used to access the Secrets Manager)
- "AWS_REGION" (The region where the Secrets Manager instance is defined)


### GCP
- "GOOGLE_APPLICATION_CREDENTIALS" (Contains the path to the json file with the credentials for your google cloud account)
- "GCP_parent" (Contains the reference to the parent of the secrets in the format 'projects/parentid')

For AWS the following permissions are needed:
  - secretsmanager:Name
  - secretsmanager:Description
  - secretsmanager:KmsKeyId
  - aws:RequestTag/${TagKey}
  - aws:ResourceTag/${TagKey}
  - aws:TagKeys
  - secretsmanager:ResourceTag/tag-key
  - secretsmanager:AddReplicaRegions
  - secretsmanager:ForceOverwriteReplicaSecret
  - resourcetypes: Secret*
  - secretsmanager:TagResource
  - secretsmanager:UntagResource

## Path specification

The path is a logical/human understandable approach to remember where the secrets are stored, the path can be any of the following formats:
- /some/words/to/define
- some/words/to/path
- /a/path/somewhere/
- a/path/to/some/secret

- Any "/" at the beginning or end of the path will be removed


## Write secrets

This command will write a secret into the secret store using the path specified and the name and value.

harpo <cloud_provider> -write -path <path> -name <secret_name> -value <secret_value>

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev
- secret_name: This is the name the secret will have, it can be any alphanumeric with a maximum of 20
  examples:
    - servername
    - serverport
    - connectionstring
- secret_value: Will contain the value you want to store for this secret
  examples:
    - myserver.com
    - 8080
    - database1.server.com:4333

### Read secrets (Human readable)

This command will look for the secret with <secret_name> in the path <path> and will output the value in a human readable format.
This is useful when you are looking for a value in the secret store. For automations check the "Export" command.

harpo <cloud_provider> -read -path <path> -name <secret_name>

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev
- secret_name: This is the name the secret will have, it can be any alphanumeric with a maximum of 20
  examples:
    - servername
    - serverport
    - connectionstring

Example output:
```
harpo read -path /infra/dev -name serverport
The value of the secret is:  443
```

### Export secrets (Automation)

This command will output the value of a <secret_name> found in the path <path> and will output the value without formatting.
This is the best option for automation.

harpo <cloud_provider> -export -path <path> -name <secret_name>

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev
- secret_name: This is the name the secret will have, it can be any alphanumeric with a maximum of 20
  examples:
    - servername
    - serverport
    - connectionstring

Example output:
```
harpo export -path /infra/dev -name serverport
443
```

### List secrets 

This command is useful to look for a secret in a path when you don't know the secrets stored in a particular path.

harpo <cloud_provider> -list -path <path>

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev

Example output:
```
harpo list -path /infra/dev
The path for the secret is:  infra/dev/serverport
The path for the secret is:  infra/dev/servername
```

### Delete secrets

This command will delete a secret from a cloud provider, `there is no confirmation requested`. Keep in mind each cloud provider has a retention policy configuration, by default when you delete a secret they stay "hidden" for certain time which makes the creation of another secret with the same name impossible until that grace period is expired, please check your cloud provider documentation for more information.

harpo <cloud_provider> delete -path <path> -name <secret_name> 

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev
- secret_name: This is the name the secret will have, it can be any alphanumeric with a maximum of 20
  examples:
    - servername
    - serverport
    - connectionstring

Example output:
```
harpo az delete -path /infra/demo -name servername
Deleteing secret from Azure Key Vault
Successfully deleted secret from Azure Keyvault
```
### Copy secrets

This command will copy a secret from one cloud provider to another one, is only possible to copy from one cloud provider to another.

harpo <cloud_provider> copy -path <path> -name <secret_name> -destination <cloud_provider>

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev
- secret_name: This is the name the secret will have, it can be any alphanumeric with a maximum of 20
  examples:
    - servername
    - serverport
    - connectionstring

Example output:
```
harpo aws copy -path /infra/dev/servers/internal/sql/primary/ -name sqlconnectionstring -destination az
Copying secret from AWS to AZ
```

### Migrate secrets

This command will copy the secret from one cloud provider to another and then delete the origin one thus moving the secret.

harpo <cloud_provider> migrate -path <path> -name <secret_name> -destination <cloud_provider>

- cloud_provider: Could be "az" for Azure or "aws" for AWS
- path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev
- secret_name: This is the name the secret will have, it can be any alphanumeric with a maximum of 20
  examples:
    - servername
    - serverport
    - connectionstring

Example output:
```
harpo aws migrate -path /infra/dev/servers/internal/sql/primary/ -name sqlconnectionstring -destination az
Migrating secret from AWS to AZ
```
