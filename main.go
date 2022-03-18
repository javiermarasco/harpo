package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	flags := os.Args
	if len(flags) < 2 {
		fmt.Println("Usage: 'kvcli write/read' ")
		return
	}
	provider := os.Args[1]
	if provider != "az" && provider != "aws" {
		fmt.Println("Provider should be 'az' or 'aws'")
		return
	}
	provider = strings.ToUpper(provider)
	operation := os.Args[2]
	operation = strings.ToUpper(operation)

	path := os.Args[4]
	var name string
	var value string
	if len(os.Args) >= 6 {
		name = os.Args[6]
	}
	if len(os.Args) >= 9 {
		value = os.Args[8]
	}

	switch operation {
	case "WRITE":
		switch provider {
		case "AZ":
			if len(os.Args) < 8 {
				fmt.Println("A path, name and value needs to be provided.")
				fmt.Println("Example: secretscli.exe (provider) write -path /infra/dev -name portnumber -value 8080 ")
				return
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(path)
			GetToken(&credentials)

			var secret secret_struct
			secret.Name = name
			secret.Value = value

			WriteAzSecret(newpath, secret, &credentials)

		case "AWS":
			newpath := SanitizePath(path)
			WriteAWSSecret(newpath, "eu-central-1", name, value)

		default:
		}

	case "READ":
		switch provider {
		case "AZ":
			if len(os.Args) < 6 {
				fmt.Println("A path and name needs to be provided.")
				fmt.Println("Example: secretscli.exe (provider) read -path /infra/dev -name portnumber ")
				return
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(path)
			GetToken(&credentials)

			var secret secret_struct
			secret.Name = name

			value, err := ReadAzSecret(newpath, secret, &credentials)
			if err != "" {
				fmt.Println(err)
				return
			}
			fmt.Println("The value of the secret is: ", value)

		case "AWS":
			newpath := SanitizePath(path)
			value, err := ReadAWSSecret(newpath, name)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from AWS.")
				return
			}
			fmt.Println("The value of the secret is: ", value)
		default:

		}

	case "EXPORT":
		switch provider {
		case "AZ":
			if len(os.Args) < 6 {
				fmt.Println("A path and name needs to be provided.")
				fmt.Println("Example: secretscli.exe (provider) export -path /infra/dev -name portnumber ")
				return
			}
			credentials := GetAzureConfig()
			// name =>  KV_<secretname>
			newpath := SanitizePath(path)
			GetToken(&credentials)

			var secret secret_struct
			secret.Name = name

			value, err := ReadAzSecret(newpath, secret, &credentials)
			if err != "" {
				fmt.Println(err)
				return
			}
			fmt.Println(value)
		case "AWS":
			newpath := SanitizePath(path)
			value, err := ReadAWSSecret(newpath, name)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from AWS.")
				return
			}
			fmt.Println(value)
		default:
		}

	case "LIST":

		switch provider {

		case "AZ":
			if len(os.Args) < 4 {
				fmt.Println("A path needs to be provided.")
				fmt.Println("Example: secretscli.exe (provider) list -path /infra/dev ")
				return
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(path)
			GetToken(&credentials)

			listOfSecrets := GetAzSecrets(&credentials)
			for _, item := range listOfSecrets.Value {
				current_path := TagsToPath(item.Tags, newpath)
				if current_path != "" {
					fmt.Println("The path for the secret is: ", current_path)
				}

			}
		case "AWS":
			//newpath := SanitizePath(path)
		default:
		}

	default:
		fmt.Println("Usage: 'kvcli write/read/export/list' ")
	}

}
