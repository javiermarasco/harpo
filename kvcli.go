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

type auth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Resource     string `json:"resource"`
	TenantID     string `json:"tenant_id"`
	Token        string `json:"access_token"`
}

func main() {
	Client_Id := os.Getenv("KV_clientid")
	Client_Secret := os.Getenv("KV_clientsecret")
	Tenant_Id := os.Getenv("KV_tenantid")
	KeyvaultName := os.Getenv("KV_kvname")

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

	// var path string
	// myPath := flag.NewFlagSet("path", flag.ExitOnError)
	// myPath.StringVar(&path, "path", "", "A path to a secret")
	// myPath.Parse(os.Args[2:])

	// var name string
	// myName := flag.NewFlagSet("name", flag.ExitOnError)
	// myName.StringVar(&name, "name", "", "A secret name")
	// myName.Parse(os.Args[3:])

	// var value string
	// myValue := flag.NewFlagSet("value", flag.ExitOnError)
	// myValue.StringVar(&value, "value", "", "A secret value")
	// myValue.Parse(os.Args[4:])

	flag.Parse()

	flagArgs := flag.Args()

	switch operation {
	case "WRITE":
		fmt.Println("Write")
		if len(flagArgs) < 7 {
			fmt.Println("A path, name and value needs to be provided.")
			fmt.Println("Example: kvcli.exe write -path /infra/dev -name portnumber -value 8080 ")
			return
		}
		fmt.Println("Path es: ", flagArgs[2][0:1])
		// if flagArgs[0] == "" {
		// 	fmt.Println("A path needs to be specified.")
		// 	return
		// }
		// if flagArgs[1] == "" {
		// 	fmt.Println("A name needs to be specified.")
		// 	return
		// }
		// if flagArgs[2] == "" {
		// 	fmt.Println("A value needs to be specified.")
		// 	return
		// }
		newpath := ""
		if flagArgs[2][0:1] == "/" {
			runes := []rune(flagArgs[2])
			newpath = string(runes[1:len(flagArgs[2])])
		}

		GetToken2(&credentials)

		var secret secret_struct
		secret.Name = flagArgs[4]
		secret.Value = flagArgs[6]

		WriteSecret(KeyvaultName, newpath, secret, &credentials)

		//os.Args[2] Needs to be the path
		//os.Args[3] Needs to be the name
		//os.Args[4] Needs to be the value

		// path needs to be alphanumeric and "/" only.

		//keyvaultname = KeyvaultName

		// path := "/infra/port/dev/servilleta"
		// if path[0:1] == "/" {
		// 	runes := []rune(path)
		// 	path = string(runes[1:len(path)])
		// }
		// GetToken2(&credentials)

		// var secreto secret_struct
		// secreto.Name = "nuevosecreto_servilleta"
		// secreto.Value = "supersecreto_servilleta"

		//WriteSecret(keyvaultname, path, secreto, &credentials)

	case "READ":
		fmt.Println("Read")
		if path == "" {
			fmt.Println("A path needs to be specified.")
			return
		}
		if name == "" {
			fmt.Println("A name needs to be specified.")
			return
		}
		if path[0:1] == "/" {
			runes := []rune(path)
			path = string(runes[1:len(path)])
		}

		GetToken2(&credentials)
		var secret secret_struct
		secret.Name = name

		ReadSecret(KeyvaultName, path, secret, &credentials)

		//os.Args[2] Needs to be the path
		//os.Args[3] Needs to be the name

		//keyvaultname = KeyvaultName

		// path := "/infra/port/dev/servilleta"
		// if path[0:1] == "/" {
		// 	runes := []rune(path)
		// 	path = string(runes[1:len(path)])
		// }
		// GetToken2(&credentials)

		// var secreto secret_struct
		// secreto.Name = "nuevosecreto_servilleta"
		// secreto.Value = "supersecreto_servilleta"
		//ReadSecret(keyvaultname, path, secreto, &credentials)

	default:
		fmt.Println("Usage: 'kvcli write/read' ")
	}

}

// funciones:
//
// Get token de azure SPN
// Crear hash para nombre de secret
// Split path y devolver array de tags para poner en el secret
// Escribir secret a KV
// Leer secret de KV
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func GetToken2(credentials *auth) {
	param := url.Values{}
	param.Add("grant_type", "client_credentials")
	param.Add("client_id", credentials.ClientID)
	param.Add("client_secret", credentials.ClientSecret)
	param.Add("resource", credentials.Resource)

	uri := fmt.Sprint("https://login.microsoftonline.com/", credentials.TenantID, "/oauth2/token")

	req, err := http.NewRequest("POST", uri, strings.NewReader(param.Encode()))
	if err != nil {
		//Failed to read response.
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

		//fmt.Println(PrettyPrint(response))
		//fmt.Println(response.AccessToken)
		credentials.Token = response.AccessToken

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp)
	}

}

func GetToken(client_id string, secret string, tenant_id string) {
	param := url.Values{}
	param.Add("grant_type", "client_credentials")
	param.Add("client_id", client_id)
	param.Add("client_secret", secret)
	param.Add("resource", "https://management.azure.com/")

	uri := fmt.Sprint("https://login.microsoftonline.com/", tenant_id, "/oauth2/token")

	req, err := http.NewRequest("POST", uri, strings.NewReader(param.Encode()))
	if err != nil {
		//Failed to read response.
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

		//fmt.Println(PrettyPrint(response))
		fmt.Println(response.AccessToken)

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

func TagsToPath(path *string, tags map[string]string) {
	// Reconstruimos un path en base a los tags?
	// Necesitamos esto?
}

func WriteSecret(keyvault_name string, path string, secret secret_struct, creds *auth) {
	base_uri := fmt.Sprint("https://", keyvault_name, ".vault.azure.net")

	////////////////////////////////////
	// Agarrar el path y agregarle "+secret.name" DONE
	// A eso hacerle un hash con CreateHash DONE
	// El hash lo ponemos como nombre de secret en la linea de uri := BLAH DONE
	// Agregamos un tag llamado "SecretName" y de valor le ponemos lo que teniamos en secret.name
	////////////////////////////////////

	// Add the SecretName as part of the path to make the hash that will be the name in Azure
	inputForHash := path + "+" + secret.Name
	secretName := CreateHash(inputForHash)

	uri := fmt.Sprint(base_uri, "/secrets/", secretName, "?api-version=7.2")

	var newtags = make(map[string]string)
	newtags = PathToTags(path)

	// Add secretname as tag
	newtags["SecretName"] = secret.Name

	type Secret struct {
		// Application specific metadata in the form of key-value pairs.
		Tags map[string]string `json:"tags,omitempty"`
		// The secret value.
		Value string `json:"value,omitempty"`
	}

	secreto := Secret{
		Value: secret.Value,
		Tags:  newtags,
	}

	jsonData, err := json.Marshal(secreto)
	// PUT
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Host","management.azure.com")
	if err != nil {
		//Failed to read response.
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
			//Failed to read response.
			panic(err)
		}
		var response secret_struct

		if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		fmt.Println(PrettyPrint(response))
		//fmt.Println(response.AccessToken)

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp)
	}

}

func ReadSecret(keyvault_name string, path string, secret_name secret_struct, creds *auth) {
	// Mismo funcionamiento que el write pero a la inversa,
	// Del path obtenemos los tags
	// Obtenemos el nombre del secret de secret.Name
	// Buscamos un secret en "keyvault_name" en base a los tags?
	// tener en cuenta el orden de los tags y buscarlos en ese orden para respetar la estructura de folders.

	// GET {vaultBaseUrl}/secrets/{secret-name}/{secret-version}?api-version=7.2
	secretname := CreateHash(path + "+" + secret_name.Name)
	base_uri := fmt.Sprint("https://", keyvault_name, ".vault.azure.net")
	uri := fmt.Sprint(base_uri, "/secrets/", secretname, "?api-version=7.2")

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
			//Failed to read response.
			panic(err)
		}
		var response secret_struct

		if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		//fmt.Println(PrettyPrint(response))
		fmt.Println("Secret value is: ", response.Value)

	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Secret not found, check the path or the secret name.")
	}

}
