package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAzSecrets(creds *auth) secret_list {

	base_uri := fmt.Sprint("https://", creds.KeyVault, ".vault.azure.net")
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
