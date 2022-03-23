package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	flags := os.Args
	if len(flags) < 2 {
		fmt.Println("Usage: 'kvcli write/read' ")
		return
	}
	provider := os.Args[1]
	if provider != "az" && provider != "aws" && provider != "gcp" {
		fmt.Println("Provider should be 'az', 'aws' or 'gcp'")
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
			WriteAWSSecret(newpath, name, value)

		case "GCP":
			newpath := SanitizePath(path)
			WriteGCPSecret(newpath, name, value)
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

		case "GCP":
			newpath := SanitizePath(path)
			value, err := ReadGCPSecret(newpath, name)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from GCP.")
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
		case "GCP":
			newpath := SanitizePath(path)
			value, err := ReadGCPSecret(newpath, name)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from GCP.")
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
			newpath := SanitizePath(path)

			secrets, err := GetAwsSecrets(newpath)
			if err != nil {
				fmt.Println("An error ocurred while retrieving the secret from AWS.")
				return
			}

			for _, secret := range secrets {
				var secretpath string
				var secretname string
				for _, value := range secret.Tags {

					if aws.ToString(value.Key) == "Path" {
						secretpath = aws.ToString(value.Value)

					}
					if aws.ToString(value.Key) == "SecretName" {
						secretname = aws.ToString(value.Value)
					}

				}
				fmt.Println("The path for the secret is: ", secretpath+"/"+secretname)
			}
		case "GCP":
			newpath := SanitizePath(path)
			GetGcpSecrets(newpath)

		default:
		}

	case "DELETE":
		switch provider {
		case "AZ":
			newpath := SanitizePath(path)
			credentials := GetAzureConfig()
			GetToken(&credentials)
			fmt.Println("Deleteing secret from Azure Key Vault")
			DeleteAzSecret(newpath, secret_struct{Name: name}, &credentials)
		case "AWS":
			newpath := SanitizePath(path)
			fmt.Println("Deleteing secret from AWS secrets manager")
			DeleteAwsSecret(newpath, name)
		case "GCP":
			fmt.Println("Deleteing secret from GCP secret manager")
			newpath := SanitizePath(path)
			DeleteGCPSecret(newpath, name)
		}
	case "COPY":
		switch provider {
		case "AWS":
			destination := strings.ToUpper(os.Args[8])
			newpath := SanitizePath(path)
			//Read from AWS
			secretvalue, err := ReadAWSSecret(newpath, name)
			if err != nil {
				fmt.Println("An error found while reading the AWS secret")
			}
			switch destination {
			case "AZ":
				fmt.Println("Copying secret from AWS to AZ")
				credentials := GetAzureConfig()
				GetToken(&credentials)
				// Write to AZ
				WriteAzSecret(newpath, secret_struct{Name: name, Value: secretvalue}, &credentials)
			case "AWS":
				fmt.Println("Not possible to copy from AWS to AWS")

			case "GCP":
				// write to GCP
				fmt.Println("Copying secret from AWS to GCP")
				WriteGCPSecret(newpath, name, secretvalue)
			default:
			}
		case "AZ":
			// copy from AZ to AWS
			destination := strings.ToUpper(os.Args[8])
			newpath := SanitizePath(path)
			//read from AZ
			credentials := GetAzureConfig()
			GetToken(&credentials)
			value, err := ReadAzSecret(newpath, secret_struct{Name: name}, &credentials)
			if err != "" {
				fmt.Println(err)
				return
			}

			switch destination {
			case "AZ":
				fmt.Println("Not possible to copy from AZ to AZ")
			case "AWS":
				// Write to AWS
				fmt.Println("Copying secret from AZ to AWS")
				WriteAWSSecret(newpath, name, value)
			case "GCP":
				// write to GCP
				fmt.Println("Copying secret from AZ to GCP")
				WriteGCPSecret(newpath, name, value)

			default:

			}
		case "GCP":
			// copy from GCP to AZ
			destination := strings.ToUpper(os.Args[8])
			newpath := SanitizePath(path)

			// Read from GCP
			value, err := ReadGCPSecret(newpath, name)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from GCP.")
				return
			}
			switch destination {
			case "AZ":
				// Write in AZ
				fmt.Println("Copying secret from GCP to AZ")
				credentials := GetAzureConfig()
				GetToken(&credentials)
				// Write to AZ
				WriteAzSecret(newpath, secret_struct{Name: name, Value: value}, &credentials)
			case "AWS":
				// Write in AWS
				fmt.Println("Copying secret from GCP to AWS")
				WriteAWSSecret(newpath, name, value)

			case "GCP":
				fmt.Println("Not possible to copy from GCP to GCP")
			default:
			}
		case "MIGRATE":
			switch provider {
			case "AWS":
				destination := strings.ToUpper(os.Args[8])
				newpath := SanitizePath(path)
				//Read from AWS
				secretvalue, err := ReadAWSSecret(newpath, name)
				if err != nil {
					fmt.Println("An error found while reading the AWS secret")
				}
				switch destination {
				case "AZ":
					fmt.Println("Migrating secret from AWS to AZ")
					credentials := GetAzureConfig()
					GetToken(&credentials)
					// Write to AZ
					err := WriteAzSecret(newpath, secret_struct{Name: name, Value: secretvalue}, &credentials)
					if err != nil {
						fmt.Print("An error occurred while migrating the secret from AWS to AZ, the origin will not be deleted.")
					} else {
						fmt.Print("Successfully migrated secret from AWS to AZ")
						DeleteAwsSecret(newpath, name)
					}

					DeleteAwsSecret(newpath, name)
				case "AWS":
					fmt.Println("Not possible to migrate from AWS to AWS")

				case "GCP":
					// write to GCP
					fmt.Println("Migrating secret from AWS to GCP")
					err := WriteGCPSecret(newpath, name, secretvalue)
					if err != nil {
						fmt.Print("An error occurred while migrating the secret from AWS to GCP, the origin will not be deleted.")
					} else {
						fmt.Print("Successfully migrated secret from AWS to GCP")
						DeleteAwsSecret(newpath, name)
					}

				default:
				}
			case "AZ":
				// copy from AZ to AWS
				destination := strings.ToUpper(os.Args[8])
				newpath := SanitizePath(path)
				//read from AZ
				credentials := GetAzureConfig()
				GetToken(&credentials)
				value, err := ReadAzSecret(newpath, secret_struct{Name: name}, &credentials)
				if err != "" {
					fmt.Println(err)
					return
				}

				switch destination {
				case "AZ":
					fmt.Println("Not possible to migrate from AZ to AZ")
				case "AWS":
					// Write to AWS
					fmt.Println("Migrating secret from AZ to AWS")
					err := WriteAWSSecret(newpath, name, value)
					if err != nil {
						fmt.Print("An error occurred while migrating the secret from AZ to AWS, the origin will not be deleted.")
					} else {
						fmt.Print("Successfully migrated secret from AZ to AWS")
						credentials := GetAzureConfig()
						GetToken(&credentials)
						DeleteAzSecret(newpath, secret_struct{Name: name}, &credentials)
					}

				case "GCP":
					// write to GCP
					fmt.Println("Migrating secret from AZ to GCP")
					err := WriteGCPSecret(newpath, name, value)
					if err != nil {
						fmt.Print("An error occurred while migrating the secret from AZ to GCP, the origin will not be deleted.")
					} else {
						fmt.Print("Successfully migrated secret from AZ to GCP")
						credentials := GetAzureConfig()
						GetToken(&credentials)
						DeleteAzSecret(newpath, secret_struct{Name: name}, &credentials)
					}

				default:

				}
			case "GCP":
				// copy from GCP to AZ
				destination := strings.ToUpper(os.Args[8])
				newpath := SanitizePath(path)

				// Read from GCP
				value, err := ReadGCPSecret(newpath, name)
				if err != nil {
					fmt.Println("An error ocurred while reading the secret from GCP.")
					return
				}
				switch destination {
				case "AZ":
					fmt.Println("Migrating secret from GCP to AZ")
					credentials := GetAzureConfig()
					GetToken(&credentials)
					// Write to AZ
					err := WriteAzSecret(newpath, secret_struct{Name: name, Value: value}, &credentials)
					if err != nil {
						fmt.Print("An error occurred while migrating the secret from GCP to AZ, the origin will not be deleted.")
					} else {
						fmt.Print("Successfully migrated secret from GCP to AZ")
						DeleteGCPSecret(newpath, name)
					}

					DeleteAwsSecret(newpath, name)
				case "AWS":
					// Write to AWS
					fmt.Println("Migrating secret from GCP to AWS")
					err := WriteAWSSecret(newpath, name, value)
					if err != nil {
						fmt.Print("An error occurred while migrating the secret from GCP to AWS, the origin will not be deleted.")
					} else {
						fmt.Print("Successfully migrated secret from GCP to AWS")
						DeleteGCPSecret(newpath, name)
					}

				case "GCP":
					fmt.Println("Not possible to migrate from GCP to GCP")
				default:
				}
			default:
				fmt.Println("Usage: 'kvcli write/read/export/list' ")
			}
		}
	}
}
