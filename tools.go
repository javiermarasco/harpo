package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
)

func TagsToPath(Tags map[string]string, Path string) string {
	var letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	var returnPath string
	var tagCount int = len(Tags) - 1

	var inputTags []string = strings.Split(Path, "/")
	var returnValue bool = false

	for index, value := range letters {
		if index < tagCount {
			returnPath += Tags[value] + "/"
		} else {
			// Split the return path just created after removing the trailing slash
			var returnPathSlice []string = strings.Split(returnPath[0:len(returnPath)-1], "/")

			// Check the positional tags in inputTags matches the ones in the just created path
			// and toggle the returnValue switch based on match
			if len(inputTags) > len(returnPathSlice) {
				returnValue = false
			} else {
				for index, _ := range inputTags {
					if returnPathSlice[index] == inputTags[index] {
						returnValue = true
					} else {
						returnValue = false
					}
				}
			}

			break
		}
	}
	if returnValue {
		returnPath += Tags["SecretName"]
		return returnPath
	} else {
		return ""
	}
}

func PathToTags(path string) map[string]string {
	pathElements := strings.Split(path, "/")
	tags := make(map[string]string)
	for i := 0; i < len(pathElements); i++ {
		index := string('A' + i)
		tags[index] = pathElements[i]
	}
	return tags
}

func SanitizePath(path string) string {
	// First clear initial slash if exist
	strings.Replace(path, "\\", "/", -1)
	var newpath string = path
	if path[0:1] == "/" {
		runes := []rune(path)
		newpath = string(runes[1:len(path)])
	}
	// Second clear trailing slash if exist
	if newpath[len(newpath)-1:len(newpath)] == "/" {
		runes := []rune(newpath)
		newpath = string(runes[0 : len(newpath)-1])
	}
	// Return sanitized path
	return newpath
}

func JsonPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func CreateHash(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func PathToGCP(path string) string {
	return strings.ToLower(strings.Replace(path, "/", "_", -1))
}

func PathFromGCP(path string) string {
	return strings.Replace(path, "_", "/", -1)
}

func DisplayHelpUsage() {
	fmt.Println("Harpo - a tool to manage secrets in Azure, AWS and GCP")
	fmt.Println()
	fmt.Println("Usage: harpo [az|aws|gcp] [action] [options]")
	fmt.Println()
	fmt.Println()
	fmt.Println("az:  Actions to interact with Azure Keyvault")
	fmt.Println("aws: Actions to interact with AWS Secrets Manager")
	fmt.Println("gcp: Actions to interact with GCP Secret Manager")
	fmt.Println()
	fmt.Println("actions:")
	fmt.Println()
	fmt.Println("List:     Will present all the secrets in a path defined by '-path'")
	fmt.Println("Write:    Will write a secret in '-path' with name '-name' with the value '-value'")
	fmt.Println("          Write will fail if the secret already exists, to update the secret use action 'Update'")
	fmt.Println("Read:     Will read the value of a secret in '-path' with name '-name'")
	fmt.Println("Update:   Will update a secret in path '-path' that has a name '-name' with the value defined in '-updatedvalue'")
	fmt.Println("Delete:   Will delete a secret in path '-path' with name '-name'")
	fmt.Println("Copy:     Will copy a secret from path '-path' with name '-name' from the cloud defined in the action [az|aws|gcp] to the cloud defined in '-destination'")
	fmt.Println("Migrate:  Will migrate the secret from the path '-path' with name '-name' to cloud defined in '-destination'.")
	fmt.Println("          After copy the secret, it will be deleted from the origin cloud provider")
	fmt.Println()

}
