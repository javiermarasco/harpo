package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type secret_struct struct {
	Name  string            `json:"name"`
	Value string            `json:"value"`
	Tags  map[string]string `json:"tags"`
}

type secret_list struct {
	Value []struct {
		ID         string `json:"id"`
		Attributes struct {
			Enabled         bool   `json:"enabled"`
			Created         int    `json:"created"`
			Updated         int    `json:"updated"`
			RecoveryLevel   string `json:"recoveryLevel"`
			RecoverableDays int    `json:"recoverableDays"`
		} `json:"attributes"`
		Tags map[string]string `json:"tags"`
	} `json:"value"`
	NextLink string `json:"nextLink"`
}
type auth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Resource     string `json:"resource"`
	TenantID     string `json:"tenant_id"`
	Token        string `json:"access_token"`
}

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

func SanitizePath(path string) string {
	// First clear initial slash if exist
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

func GetToken(credentials *auth) {
	param := url.Values{}
	param.Add("grant_type", "client_credentials")
	param.Add("client_id", credentials.ClientID)
	param.Add("client_secret", credentials.ClientSecret)
	param.Add("resource", credentials.Resource)

	uri := fmt.Sprint("https://login.microsoftonline.com/", credentials.TenantID, "/oauth2/token")

	req, err := http.NewRequest("POST", uri, strings.NewReader(param.Encode()))
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	type response_token struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    string `json:"expires_in"`
		ExtExpiresIn string `json:"ext_expires_in"`
		ExpiresOn    string `json:"expires_on"`
		NotBefore    string `json:"not_before"`
		Resource     string `json:"resource"`
		AccessToken  string `json:"access_token"`
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//Failed to read response.
			panic(err)
		}
		var response response_token

		if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		//fmt.Println(JsonPrint(response))
		//fmt.Println(response.AccessToken)
		credentials.Token = response.AccessToken

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp)
	}

}

func CreateHash(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
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
			for index, _ := range inputTags {
				if returnPathSlice[index] == inputTags[index] {
					returnValue = true
				} else {
					returnValue = false
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

func WriteSecret(keyvault_name string, path string, secret secret_struct, creds *auth) {
	base_uri := fmt.Sprint("https://", keyvault_name, ".vault.azure.net")
	inputForHash := path + "+" + secret.Name
	secretName := CreateHash(inputForHash)

	uri := fmt.Sprint(base_uri, "/secrets/", secretName, "?api-version=7.2")

	var newtags = make(map[string]string)
	newtags = PathToTags(path)

	newtags["SecretName"] = secret.Name

	type Secret struct {
		Tags  map[string]string `json:"tags,omitempty"`
		Value string            `json:"value,omitempty"`
	}

	secreto := Secret{
		Value: secret.Value,
		Tags:  newtags,
	}

	jsonData, _ := json.Marshal(secreto)
	// PUT
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var response secret_struct

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
	} else {
		fmt.Println("Get failed with error: ", resp)
	}

}

func ReadSecret(keyvault_name string, path string, secret_name secret_struct, creds *auth) (value string, error string) {
	secretname := CreateHash(path + "+" + secret_name.Name)
	base_uri := fmt.Sprint("https://", keyvault_name, ".vault.azure.net")
	uri := fmt.Sprint(base_uri, "/secrets/", secretname, "?api-version=7.2")
	value = ""
	error = ""
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var response secret_struct

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		value = response.Value

	} else if resp.StatusCode == http.StatusNotFound {
		error = "Secret not found, check the path or the secret name."
	}
	return value, error
}

func GetSecrets(keyvault_name string, creds *auth) secret_list {

	base_uri := fmt.Sprint("https://", keyvault_name, ".vault.azure.net")
	uri := fmt.Sprint(base_uri, "/secrets?api-version=7.2")
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	var response secret_list
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Secret not found, check the path or the secret name.")
	}
	return response
}
