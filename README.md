# Secret Manager CLI

The objective of this cli is to be able to manage secrets in Azure, AWS and GCP by using a path to 
specify where a secret is stored in a "folder structure" approach while providing an easy way to list 
secrets in this folder structure and also an easy way to retrieve the values of those secrets.

## [You can follow the progress of this project on my live streams on Twitch or Youtube](https://linktr.ee/javi__codes)

# How to use

## Azure

You need to define 4 environment variables
- "AZ_clientid" (Contains the client id of the service princial/app registration used to access your keyvault)
- "AZ_clientsecret" (Contains the secret of the service principal/app registration)
- "AZ_tenantid" (Contains the Tenant ID where your keyvault is deployed)
- "AZ_kvname" (Contains the name of the keyvault to use)

Once your have this variables defined you can start using the CLI.

### Write secrets

This command will write a secret into the secret store using the path specified and the name and value.

secretmanager -write -path <secret_path> -name <secret_name> -value <secret_value>

- secret_path: This is the path where the secret will be stored, it can start or end with a "/"
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

This command will look for the secret with <secret_name> in the path <secret_path> and will output the value in a human readable format.
This is useful when you are looking for a value in the secret store. For automations check the "Export" command.

secretmanager -read -path <secret_path> -name <secret_name>

- secret_path: This is the path where the secret will be stored, it can start or end with a "/"
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
kvcli.exe read -path /infra/dev -name serverport
The value of the secret is:  443
```

### Export secrets (Automation)

This command will output the value of a <secret_name> found in the path <secret_path> and will output the value without formatting.
This is the best option for automation.

secretmanager -export -path <secret_path> -name <secret_name>

- secret_path: This is the path where the secret will be stored, it can start or end with a "/"
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
kvcli.exe export -path /infra/dev -name serverport
443
```

### List secrets 

This command is useful to look for a secret in a path when you don't know the secrets stored in a particular path.

secretmanager -list -path <secret_path>

- secret_path: This is the path where the secret will be stored, it can start or end with a "/"
  examples: 
    - /infra/dev
    - /infra/dev/
    - infra/dev

Example output:
```
kvcli.exe list -path /infra/dev
The path for the secret is:  infra/dev/serverport
The path for the secret is:  infra/dev/servername
```

## AWS (TODO)

## GCP (TODO)