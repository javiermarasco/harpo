package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	Client_Id := os.Getenv("AZ_clientid")
	Client_Secret := os.Getenv("AZ_clientsecret")
	Tenant_Id := os.Getenv("AZ_tenantid")
	KeyvaultName := os.Getenv("AZ_kvname")

	credentials := auth{ClientID: Client_Id, ClientSecret: Client_Secret, Resource: "https://vault.azure.net", TenantID: Tenant_Id}

	flags := os.Args
	if len(flags) < 2 {
		fmt.Println("Usage: 'kvcli write/read' ")
		return
	}
	operation := os.Args[1]
	operation = strings.ToUpper(operation)

	var path string
	var name string
	var value string

	flag.StringVar(&path, "path", "", "a path to a secret")
	flag.StringVar(&name, "name", "", "a secret name")
	flag.StringVar(&value, "value", "", "a secret value")

	flag.Parse()

	flagArgs := flag.Args()

	switch operation {
	case "WRITE":
		if len(flagArgs) < 7 {
			fmt.Println("A path, name and value needs to be provided.")
			fmt.Println("Example: kvcli.exe write -path /infra/dev -name portnumber -value 8080 ")
			return
		}

		newpath := SanitizePath(flagArgs[2])
		GetToken(&credentials)

		var secret secret_struct
		secret.Name = flagArgs[4]
		secret.Value = flagArgs[6]

		WriteSecret(KeyvaultName, newpath, secret, &credentials)

	case "READ":
		if len(flagArgs) < 5 {
			fmt.Println("A path and name needs to be provided.")
			fmt.Println("Example: kvcli.exe read -path /infra/dev -name portnumber ")
			return
		}

		newpath := SanitizePath(flagArgs[2])
		GetToken(&credentials)

		var secret secret_struct
		secret.Name = flagArgs[4]

		value, err := ReadSecret(KeyvaultName, newpath, secret, &credentials)
		if err != "" {
			fmt.Println(err)
			return
		}
		fmt.Println("The value of the secret is: ", value)

	case "EXPORT":
		if len(flagArgs) < 5 {
			fmt.Println("A path and name needs to be provided.")
			fmt.Println("Example: kvcli.exe export -path /infra/dev -name portnumber ")
			return
		}

		// name =>  KV_<secretname>
		newpath := SanitizePath(flagArgs[2])
		GetToken(&credentials)

		var secret secret_struct
		secret.Name = flagArgs[4]

		value, err := ReadSecret(KeyvaultName, newpath, secret, &credentials)
		if err != "" {
			fmt.Println(err)
			return
		}
		fmt.Println(value)

	case "LIST":
		if len(flagArgs) < 3 {
			fmt.Println("A path needs to be provided.")
			fmt.Println("Example: kvcli.exe list -path /infra/dev ")
			return
		}
		newpath := SanitizePath(flagArgs[2])
		GetToken(&credentials)

		listOfSecrets := GetSecrets(KeyvaultName, &credentials)
		for _, item := range listOfSecrets.Value {
			current_path := TagsToPath(item.Tags, newpath)
			if current_path != "" {
				fmt.Println("The path for the secret is: ", current_path)
			}

		}

	default:
		fmt.Println("Usage: 'kvcli write/read/export/list' ")
	}

}
