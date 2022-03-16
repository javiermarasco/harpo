package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetAzureConfig() auth {
	Client_Id := os.Getenv("AZ_clientid")
	Client_Secret := os.Getenv("AZ_clientsecret")
	Tenant_Id := os.Getenv("AZ_tenantid")
	KeyvaultName := os.Getenv("AZ_kvname")

	return auth{ClientID: Client_Id, ClientSecret: Client_Secret, Resource: "https://vault.azure.net", TenantID: Tenant_Id, KeyVault: KeyvaultName}

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
