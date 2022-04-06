package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {

	// Define all the flag sets for the commands
	azCmd := flag.NewFlagSet("az", flag.ExitOnError)
	azPath := azCmd.String("path", " ", "This is the path for the secret")
	azName := azCmd.String("name", " ", "This is the name for the secret")
	azValue := azCmd.String("value", " ", "This is the value for the secret")
	azNewValue := azCmd.String("newvalue", " ", "This is the new value for the secret")
	azDestination := azCmd.String("destination", " ", "This is the destination cloud provider for copy and migrate a secret")

	awsCmd := flag.NewFlagSet("aws", flag.ExitOnError)
	awsPath := awsCmd.String("path", " ", "This is the path for the secret")
	awsName := awsCmd.String("name", " ", "This is the name for the secret")
	awsValue := awsCmd.String("value", " ", "This is the value for the secret")
	awsNewValue := awsCmd.String("newvalue", " ", "This is the new value for the secret")
	awsDestination := awsCmd.String("destination", " ", "This is the destination cloud provider for copy and migrate a secret")

	gcpCmd := flag.NewFlagSet("gcp", flag.ExitOnError)
	gcpPath := gcpCmd.String("path", " ", "This is the path for the secret")
	gcpName := gcpCmd.String("name", " ", "This is the name for the secret")
	gcpValue := gcpCmd.String("value", " ", "This is the value for the secret")
	gcpNewValue := gcpCmd.String("newvalue", " ", "This is the new value for the secret")
	gcpDestination := gcpCmd.String("destination", " ", "This is the destination cloud provider for copy and migrate a secret")

	// Verify at least 2 arguments are provided
	if len(os.Args) < 2 {
		DisplayHelpUsage()
		os.Exit(1)
	}

	// Verify the provided provider is valid (either az, aws or gcp)
	provider := strings.ToUpper(os.Args[1])
	if provider != "AZ" && provider != "AWS" && provider != "GCP" {
		fmt.Println("expected 'az','aws' or 'gcp'.")
		DisplayHelpUsage()
		os.Exit(1)
	}

	// Verify the provided operation is valid
	operation := strings.ToUpper(os.Args[2])
	if operation != "WRITE" && operation != "READ" && operation != "EXPORT" && operation != "LIST" && operation != "DELETE" && operation != "COPY" && operation != "MIGRATE" && operation != "UPDATE" {
		fmt.Println("Wrong operation.")
		DisplayHelpUsage()
		os.Exit(1)
	}
	switch operation {
	case "WRITE":
		switch provider {
		case "AZ":
			azCmd.Parse(os.Args[3:])
			if *azPath == " " || *azName == " " || *azValue == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(*azPath)
			GetToken(&credentials)
			var secret secret_struct
			secret.Name = *azName
			secret.Value = *azValue
			_, error := ReadAzSecret(newpath, secret, &credentials)
			if error != "" {
				WriteAzSecret(newpath, secret, &credentials)
			} else {
				fmt.Println("The specified secret already exist. Update it instead of Write.")
			}

		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " || *awsName == " " || *awsValue == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*awsPath)
			_, error := ReadAWSSecret(newpath, *awsName)
			if error != nil {
				WriteAWSSecret(newpath, *awsName, *awsValue)
			} else {
				fmt.Println("The specified secret already exist. Update it instead of Write.")
			}

		case "GCP":
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " || *gcpName == " " || *gcpValue == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*gcpPath)
			_, error := ReadGCPSecret(newpath, *gcpName)
			if error != nil {
				WriteGCPSecret(newpath, *gcpName, *gcpValue)
			} else {
				fmt.Println("The specified secret already exist. Update it instead of Write.")
			}

		default:
		}
	case "READ":
		switch provider {
		case "AZ":
			azCmd.Parse(os.Args[3:])
			if *azPath == " " || *azName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(*azPath)
			GetToken(&credentials)
			var secret secret_struct
			secret.Name = *azName
			value, err := ReadAzSecret(newpath, secret, &credentials)
			if err != "" {
				fmt.Println(err)
				return
			}
			fmt.Println("The value of the secret is: ", value)
		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " || *awsName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*awsPath)
			value, err := ReadAWSSecret(newpath, *awsName)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from AWS.")
				return
			}
			fmt.Println("The value of the secret is: ", value)

		case "GCP":
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " || *gcpName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*gcpPath)
			value, err := ReadGCPSecret(newpath, *gcpName)
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
			azCmd.Parse(os.Args[3:])
			if *azPath == " " || *azName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(*azPath)
			GetToken(&credentials)
			var secret secret_struct
			secret.Name = *azName
			value, err := ReadAzSecret(newpath, secret, &credentials)
			if err != "" {
				fmt.Println(err)
				return
			}
			fmt.Println(value)
		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " || *awsName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*awsPath)
			value, err := ReadAWSSecret(newpath, *awsName)
			if err != nil {
				fmt.Println("An error ocurred while reading the secret from AWS.")
				return
			}
			fmt.Println(value)
		case "GCP":
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " || *gcpName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*gcpPath)
			value, err := ReadGCPSecret(newpath, *gcpName)
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
			azCmd.Parse(os.Args[3:])
			if *azPath == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			credentials := GetAzureConfig()
			newpath := SanitizePath(*azPath)
			GetToken(&credentials)

			listOfSecrets := GetAzSecrets(&credentials)
			for _, item := range listOfSecrets.Value {
				current_path := TagsToPath(item.Tags, newpath)
				if current_path != "" {
					fmt.Println("The path for the secret is: ", current_path)
				}

			}
		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*awsPath)
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
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*gcpPath)
			GetGcpSecrets(newpath)

		default:
		}
	case "DELETE":
		switch provider {
		case "AZ":
			azCmd.Parse(os.Args[3:])
			if *azPath == " " || *azName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*azPath)
			credentials := GetAzureConfig()
			GetToken(&credentials)
			fmt.Println("Deleteing secret from Azure Key Vault")
			DeleteAzSecret(newpath, secret_struct{Name: *azName}, &credentials)
		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " || *awsName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*awsPath)
			fmt.Println("Deleteing secret from AWS secrets manager")
			DeleteAwsSecret(newpath, *awsName)
		case "GCP":
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " || *gcpName == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*gcpPath)
			fmt.Println("Deleteing secret from GCP secret manager")
			DeleteGCPSecret(newpath, *gcpName)
		default:
		}
	case "COPY":
		switch provider {
		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " || *awsName == " " || *awsDestination == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			newpath := SanitizePath(*awsPath)
			destination := strings.ToUpper(*awsDestination)

			//Read from AWS
			secretvalue, err := ReadAWSSecret(newpath, *awsName)
			if err != nil {
				fmt.Println("An error found while reading the AWS secret")
			}
			switch destination {
			case "AZ":
				fmt.Println("Copying secret from AWS to AZ")
				credentials := GetAzureConfig()
				GetToken(&credentials)
				// Write to AZ
				WriteAzSecret(newpath, secret_struct{Name: *awsName, Value: secretvalue}, &credentials)
			case "AWS":
				fmt.Println("Not possible to copy from AWS to AWS")

			case "GCP":
				// write to GCP
				fmt.Println("Copying secret from AWS to GCP")
				WriteGCPSecret(newpath, *awsName, secretvalue)
			default:
			}
		case "AZ":
			// copy from AZ to AWS
			azCmd.Parse(os.Args[3:])
			if *azPath == " " || *azName == " " || *azDestination == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			destination := strings.ToUpper(*azDestination)
			newpath := SanitizePath(*azPath)
			//read from AZ
			credentials := GetAzureConfig()
			GetToken(&credentials)
			value, err := ReadAzSecret(newpath, secret_struct{Name: *azName}, &credentials)
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
				WriteAWSSecret(newpath, *azName, value)
			case "GCP":
				// write to GCP
				fmt.Println("Copying secret from AZ to GCP")
				WriteGCPSecret(newpath, *azName, value)

			default:
			}
		case "GCP":
			// copy from GCP to AZ
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " || *gcpName == " " || *gcpDestination == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			destination := strings.ToUpper(*gcpDestination)
			newpath := SanitizePath(*gcpPath)

			// Read from GCP
			value, err := ReadGCPSecret(newpath, *gcpName)
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
				WriteAzSecret(newpath, secret_struct{Name: *gcpName, Value: value}, &credentials)
			case "AWS":
				// Write in AWS
				fmt.Println("Copying secret from GCP to AWS")
				WriteAWSSecret(newpath, *gcpName, value)

			case "GCP":
				fmt.Println("Not possible to copy from GCP to GCP")
			default:
			}

		}
	case "MIGRATE":
		switch provider {
		case "AWS":
			awsCmd.Parse(os.Args[3:])
			if *awsPath == " " || *awsName == " " || *awsDestination == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			destination := strings.ToUpper(*awsDestination)
			newpath := SanitizePath(*awsPath)
			//Read from AWS
			secretvalue, err := ReadAWSSecret(newpath, *awsName)
			if err != nil {
				fmt.Println("An error found while reading the AWS secret")
			}
			switch destination {
			case "AZ":
				fmt.Println("Migrating secret from AWS to AZ")
				credentials := GetAzureConfig()
				GetToken(&credentials)
				// Write to AZ
				err := WriteAzSecret(newpath, secret_struct{Name: *awsName, Value: secretvalue}, &credentials)
				if err != nil {
					fmt.Print("An error occurred while migrating the secret from AWS to AZ, the origin will not be deleted.")
				} else {
					fmt.Print("Successfully migrated secret from AWS to AZ")
					DeleteAwsSecret(newpath, *awsName)
				}

				DeleteAwsSecret(newpath, *awsName)
			case "AWS":
				fmt.Println("Not possible to migrate from AWS to AWS")

			case "GCP":
				// write to GCP
				gcpCmd.Parse(os.Args[3:])
				fmt.Println("Migrating secret from AWS to GCP")
				err := WriteGCPSecret(newpath, *gcpName, secretvalue)
				if err != nil {
					fmt.Print("An error occurred while migrating the secret from AWS to GCP, the origin will not be deleted.")
				} else {
					fmt.Print("Successfully migrated secret from AWS to GCP")
					DeleteAwsSecret(newpath, *gcpName)
				}

			default:
			}
		case "AZ":
			// copy from AZ to AWS
			azCmd.Parse(os.Args[3:])
			if *azPath == " " || *azName == " " || *azDestination == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			destination := strings.ToUpper(*azDestination)
			newpath := SanitizePath(*azPath)
			//read from AZ
			credentials := GetAzureConfig()
			GetToken(&credentials)
			value, err := ReadAzSecret(newpath, secret_struct{Name: *azName}, &credentials)
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
				err := WriteAWSSecret(newpath, *azName, value)
				if err != nil {
					fmt.Print("An error occurred while migrating the secret from AZ to AWS, the origin will not be deleted.")
				} else {
					fmt.Print("Successfully migrated secret from AZ to AWS")
					credentials := GetAzureConfig()
					GetToken(&credentials)
					DeleteAzSecret(newpath, secret_struct{Name: *azName}, &credentials)
				}

			case "GCP":
				// write to GCP
				fmt.Println("Migrating secret from AZ to GCP")
				err := WriteGCPSecret(newpath, *azName, value)
				if err != nil {
					fmt.Print("An error occurred while migrating the secret from AZ to GCP, the origin will not be deleted.")
				} else {
					fmt.Print("Successfully migrated secret from AZ to GCP")
					credentials := GetAzureConfig()
					GetToken(&credentials)
					DeleteAzSecret(newpath, secret_struct{Name: *azName}, &credentials)
				}

			default:
			}
		case "GCP":
			// copy from GCP to AZ
			gcpCmd.Parse(os.Args[3:])
			if *gcpPath == " " || *gcpName == " " || *gcpDestination == " " {
				fmt.Println("Missing path,name or value")
				DisplayHelpUsage()
				os.Exit(1)
			}
			destination := strings.ToUpper(*gcpDestination)
			newpath := SanitizePath(*gcpPath)

			// Read from GCP
			value, err := ReadGCPSecret(newpath, *gcpName)
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
				err := WriteAzSecret(newpath, secret_struct{Name: *gcpName, Value: value}, &credentials)
				if err != nil {
					fmt.Print("An error occurred while migrating the secret from GCP to AZ, the origin will not be deleted.")
				} else {
					fmt.Print("Successfully migrated secret from GCP to AZ")
					DeleteGCPSecret(newpath, *gcpName)
				}
			case "AWS":
				// Write to AWS
				fmt.Println("Migrating secret from GCP to AWS")
				err := WriteAWSSecret(newpath, *gcpName, value)
				if err != nil {
					fmt.Print("An error occurred while migrating the secret from GCP to AWS, the origin will not be deleted.")
				} else {
					fmt.Print("Successfully migrated secret from GCP to AWS")
					DeleteGCPSecret(newpath, *gcpName)
				}

			case "GCP":
				fmt.Println("Not possible to migrate from GCP to GCP")
			default:
			}
		default:
		}
	case "UPDATE":
		{
			switch provider {
			case "AZ":
				azCmd.Parse(os.Args[3:])
				if *azPath == " " || *azName == " " || *azNewValue == " " {
					fmt.Println("Missing path,name or value")
					DisplayHelpUsage()
					os.Exit(1)
				}
				credentials := GetAzureConfig()
				newpath := SanitizePath(*azPath)
				GetToken(&credentials)

				var secret secret_struct
				secret.Name = *azName
				secret.Value = *azNewValue
				WriteAzSecret(newpath, secret, &credentials)

			case "AWS":
				awsCmd.Parse(os.Args[3:])
				if *awsPath == " " || *awsName == " " || *awsNewValue == " " {
					fmt.Println("Missing path,name or value")
					DisplayHelpUsage()
					os.Exit(1)
				}
				newpath := SanitizePath(*awsPath)
				UpdateAWSSecret(newpath, *awsName, *awsNewValue)

			case "GCP":
				gcpCmd.Parse(os.Args[3:])
				if *gcpPath == " " || *gcpName == " " || *gcpNewValue == " " {
					fmt.Println("Missing path,name or value")
					DisplayHelpUsage()
					os.Exit(1)
				}
				newpath := SanitizePath(*gcpPath)

				_, error := ReadGCPSecret(newpath, *gcpName)
				if error != nil {
					UpdateGCPSecret(newpath, *gcpName, *gcpNewValue)
				} else {
					fmt.Println("The specified secret already exist. Update it instead of Write.")
				}

			default:
			}
		}
	}
}
